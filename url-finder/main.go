package main

 import (
     "log"
     "net/http"
     "fmt"
     "os"
     "io/ioutil"
     "time"
     "encoding/json"
 )
var workflow_name string
var namespace string
var finalUrl string
var wrongWorkflow int

func main() {
   //Read the URL of the "workflow" and parse the workflow name and namespace
   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      namespace = r.URL.Query().Get("namespace")
      workflow_name = r.URL.Query().Get("workflowname")
      wrongWorkflow = 0
      if len(workflow_name) == 0 {
         fmt.Println("No value entered for workflow name")
      } else{
         urlGetter()
         if wrongWorkflow == 1 {
            fmt.Println("Invalid workflow name")
         } else {
            fmt.Println(finalUrl)
            http.Redirect(w, r, finalUrl, 302)
         }
      }

   })
   log.Println("Listening on :8080...")
   err := http.ListenAndServe(":8080", nil)
   if err != nil {
      log.Fatal("ListenAndServe: ", err)
   }
}

func getToken(f string, e string) string {
	if os.Getenv(e) == "" {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println("No TOKEN_FILE found. Falling back to Environment Variable")
		}
		os.Setenv(e, string(data))
	}
	a := getValidatedEnvVar(e)
	return a
}

func getValidatedEnvVar(e string) string {
	c := os.Getenv(e)
	if os.Getenv(e) == "" {
		fmt.Printf("Error: No environment variable called %s available. Exiting.\n", e)
		os.Exit(1)
	}
	return c
}

func urlGetter() {
   //Create the JSON struct
   type ItemMetadata struct {
      Name string `json:"name,omitempty"`
      Namespace string `json:"namespace,omitempty"`
      UID string `json:"uid,omitempty"`
      CreationTimestamp string `json:"creationTimestamp,omitempty"`
   }

   type Items struct {
      Metadata ItemMetadata `json:"metadata,omitempty"`
   }

   type Workflow struct {
      Items []Items `json:"items,omitempty"`
   }

   url := os.ExpandEnv("$ARGO_URL") + "/api/v1/workflows/" + namespace + "?listOptions.fieldSelector=metadata.name=" + workflow_name
   fmt.Println(url)
   
   //Set up HTTP client object
   argoClient := http.Client{Timeout: 10 * time.Second}
   
   //Create a request object
   request, err := http.NewRequest(http.MethodGet, url, nil)
   if err != nil {
      log.Fatal(err)
   }

   //Set the Authorization Token for Argo and the content type
   request.Header.Set("Authorization", getToken(os.Getenv("TOKEN_FILE"), "ACCESS_TOKEN"))
   request.Header.Set("Content-Type", "application/json")

   //Make the request and get back a response object
   response, getErr := argoClient.Do(request)
   if getErr != nil {
      log.Fatal(getErr)
   }
   //Ensure we get data back from the response
   if response.Body != nil {
      defer response.Body.Close()
   }

   //Read the data from the response object into a variable
   responseData, readErr := ioutil.ReadAll(response.Body)
   if readErr != nil {
      log.Fatalln(readErr)
   }
   
   //Create an instance of our struct
   var workflow1 *Workflow

   //Unmarshal the JSON from the body into the structure
   jsonErr := json.Unmarshal(responseData, &workflow1)
   if jsonErr != nil {
      log.Fatalln(jsonErr)
   }

   
   if workflow1.Items == nil {
      ///Check to see if workflow has been archived instead (same steps as checking Workflows)
      fmt.Println("No Workflows Found. Checking Archived Workflows instead...")
      url := os.ExpandEnv("$ARGO_URL") + "/api/v1/archived-workflows?listOptions.fieldSelector=metadata.namespace=" + namespace + ",metadata.name=" + workflow_name
      argoClient := http.Client{Timeout: 15 * time.Second}
      request, err := http.NewRequest(http.MethodGet, url, nil)
      if err != nil {
         log.Fatal(err)
      }
      request.Header.Set("Authorization", os.ExpandEnv("$ARGO_TOKEN"))
      request.Header.Set("Content-Type", "application/json")
      response, getErr := argoClient.Do(request)
      if err != nil {
         log.Fatal(getErr)
      }
      if response.Body != nil {
         defer response.Body.Close()
      }
      responseData, readErr := ioutil.ReadAll(response.Body)
      if readErr != nil {
         log.Fatalln(readErr)
      }
      var workflow1 *Workflow
      jsonErr := json.Unmarshal(responseData, &workflow1)
      if jsonErr != nil {
         log.Fatalln(jsonErr)
      }
      if workflow1.Items == nil {
         //Workflow name was invalid: Gets redirected to error page
         fmt.Println("Your workflow name was invalid")
         wrongWorkflow = 1
      } else {
         finalUrl = os.ExpandEnv("$ARGO_URL") + "archived-workflows/" + namespace + "/" + workflow1.Items[0].Metadata.UID
      }  
   } else {
      finalUrl = os.ExpandEnv("$ARGO_URL") + "workflows/" + namespace + "/" + workflow1.Items[0].Metadata.Name
   }
}
