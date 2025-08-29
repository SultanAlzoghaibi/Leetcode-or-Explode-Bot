package monitor

import (
	"Leetcode-or-Explode-Bot/internal/shared/gRPCshared"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
	"net/http"
)

func StartMonitor() {

	ipM := "127.0.0.1:4000"
	go func() {
		var status gRPCshared.Status

		for {
			status, _ = gRPCshared.StartListen(ipM)
			fmt.Printf("isSeverOn: %v\n", status)

			if !status.IsOn {

				switch status.PodName {
				case "chrome-l":
					print(status.PodName)

				case "discordbot":
					print(status.PodName)

				case "nginx":
					print(status.PodName)

				}

				//todo: impliment Reboot logic with google API
			}
		}

	}()

	gRPCshared.SendStatusTo("", true)

}

var region = "us-central"
var project = "MonitorDeploy"

func launchCloudRun(dockerImageURL string, cloudRunName string) {
	c, err := client(region)

	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
	}

	doesServiceExist, err := serviceExists(c, region, project, cloudRunName)
	if err != nil {
		fmt.Printf("Error checking if service event exists: %v\n", err)
	}
	if !doesServiceExist {
		fmt.Printf("Could not find the service to make cloud runs: %v\n", err)
		return
	}

	svc := &run.Service{
		ApiVersion: "serving.knative.dev/v1",
		Kind:       "Service",
		Metadata: &run.ObjectMeta{
			Name: cloudRunName,
		},
		Spec: &run.ServiceSpec{
			Template: &run.RevisionTemplate{
				Metadata: &run.ObjectMeta{Name: cloudRunName + "-v1"},
				Spec: &run.RevisionSpec{
					Containers: []*run.Container{
						{
							Image: "gcr.io/google-samples/hello-app:1.0",
						},
					},
				},
			},
		},
	}

	_, err = c.Namespaces.Services.Create("namespaces/"+project, svc).Do()

}

func client(region string) (*run.APIService, error) {
	return run.NewService(context.TODO(),
		option.WithEndpoint(fmt.Sprintf("https://%s-run.googleapis.com", region)))
}

func serviceExists(c *run.APIService, region, project, name string) (bool, error) {
	_, err := c.Namespaces.Services.Get(fmt.Sprintf("namespaces/%s/services/%s", project, name)).Do()
	if err == nil {
		return true, nil
	}
	// not all errors indicate service does not exist, look for 404 status code
	v, ok := err.(*googleapi.Error)
	if !ok {
		return false, fmt.Errorf("failed to query service: %w", err)
	}
	if v.Code == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code=%d from get service call: %w", v.Code, err)
}

/* CLIENC
conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := monitor.NewMonitorServiceClient(conn)
res, _ := client.SendStatus(context.Background(), &monitor.Status{ServerStatus: false})
fmt.Println("Response:", res.ServerStatus)
*/
