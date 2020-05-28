PROJECT_NO=[YOUR_PROJECT_NO]
PROJECT_ID=[YOUR_PROJECT_ID]
ENDPOINT=[YOU_ENDPOINT]
SERVICE_NAME=YOUR_SERVICE_NAME

# 初回のみ
create-service-account:
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
         --member=serviceAccount:service-${PROJECT_NO}@gcp-sa-pubsub.iam.gserviceaccount.com \
         --role=roles/iam.serviceAccountTokenCreator
	gcloud iam service-accounts create cloud-run-pubsub-invoker \
         --display-name "Cloud Run Pub/Sub Invoker"
	gcloud iam service-accounts create cloud-run-cloud-tasks-invoker \
         --display-name "Cloud Run CloudTasks Invoker"
	gcloud iam service-accounts create cloud-run-scheduler-invoker \
         --display-name "Cloud Run CloudTasks Invoker"
	gcloud run services add-iam-policy-binding ${SERVICE_NAME} \
        --member=serviceAccount:cloud-run-scheduler-invoker@${PROJECT_ID}.iam.gserviceaccount.com \
	    --role=roles/run.invoker --platform managed

# CloudRun　デプロイ
deploy:
	gcloud builds submit --tag gcr.io/${PROJECT_ID}/pubsub-practice-container
	gcloud run deploy ${SERVICE_NAME} --image gcr.io/${PROJECT_ID}/pubsub-practice-container \
	       --platform managed

# CloudRun初回時のみ付与する必要がある
setup-subscription:
	gcloud run services add-iam-policy-binding ${SERVICE_NAME} \
        --member=serviceAccount:cloud-run-pubsub-invoker@${PROJECT_ID}.iam.gserviceaccount.com \
        --role=roles/run.invoker --platform managed

	gcloud pubsub subscriptions create resizeSubscription --topic resizeTopic \
       --push-endpoint=${ENDPOINT} \
       --push-auth-service-account=cloud-run-pubsub-invoker@${PROJECT_ID}.iam.gserviceaccount.com

# CloudFunctions　デプロイ
deploy-functions:
	cd functions/image
	go mod vendor
	gcloud functions deploy ResizeSubscriber --runtime go113 --trigger-topic resizeTopic

# CloudTaskの設定-Queueの作成まで
setup-cloudtasks:
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
         --member=serviceAccount:cloud-tasks@${PROJECT_ID}.iam.gserviceaccount.com \
         --role=roles/iam.serviceAccountUser
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
        --member=serviceAccount:cloud-tasks@${PROJECT_ID}.iam.gserviceaccount.com \
        --role=roles/cloudtasks.enqueuer
	gcloud run services add-iam-policy-binding ${SERVICE_NAME} \
        --member=serviceAccount:cloud-run-cloud-tasks-invoker@${PROJECT_ID}.iam.gserviceaccount.com \
        --role=roles/run.invoker --platform managed
	gcloud tasks queues create task-queue --max-concurrent-dispatches 10 --max-attempts 1

# Schedulerの作成
setup-scheduler:
	gcloud scheduler jobs create http task-job --schedule="*/5 * * * *" \
       --http-method=POST \
       --uri=${ENDPOINT}/task \
       --oidc-service-account-email=cloud-run-scheduler-invoker@${PROJECT_ID}.iam.gserviceaccount.com   \
       --oidc-token-audience=${ENDPOINT}