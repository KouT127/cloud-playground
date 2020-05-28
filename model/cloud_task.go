package model

import (
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"context"
	"fmt"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"os"
)

type CloudTask struct{}

func CreateHTTPTaskWithToken(projectID, locationID, queueID, url, email, message string) (*taskspb.Task, error) {
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("PUBSUB_SERVICE"))))
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url + "/task",
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: email,
						},
					},
				},
			},
		},
	}

	req.Task.GetHttpRequest().Body = []byte(message)
	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}
	return createdTask, nil
}
