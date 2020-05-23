PROJECT_NO=[YOUR_PROJECT_NO]
PROJECT_ID=[YOUR_PROJECT_ID]
ENDPOINT=[YOU_ENDPOINT]

# 初回のみ
add-iam:
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
         --member=serviceAccount:service-${PROJECT_NO}@gcp-sa-pubsub.iam.gserviceaccount.com \
         --role=roles/iam.serviceAccountTokenCreator
	gcloud iam service-accounts create cloud-run-pubsub-invoker \
         --display-name "Cloud Run Pub/Sub Invoker"

# CloudRun　デプロイ
deploy:
	gcloud builds submit --tag gcr.io/${PROJECT_ID}/pubsub-practice-container
	gcloud run deploy pubsub-practice --image gcr.io/${PROJECT_ID}/pubsub-practice-container \
	       --platform managed

# CloudRun初回時のみ付与する必要がある
setup-subscription:
	gcloud run services add-iam-policy-binding pubsub-practice \
        --member=serviceAccount:cloud-run-pubsub-invoker@${PROJECT_ID}.iam.gserviceaccount.com \
        --role=roles/run.invoker
	gcloud pubsub subscriptions create resizeSubscription --topic resizeTopic \
       --push-endpoint=${ENDPOINT} \
       --push-auth-service-account=cloud-run-pubsub-invoker@${PROJECT_ID}.iam.gserviceaccount.com

# CloudFunctions　デプロイ
deploy-functions:
	cd functions/image
	go mod vendor
	gcloud functions deploy ResizeSubscriber --runtime go113 --trigger-topic resizeTopico