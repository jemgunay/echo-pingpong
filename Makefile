.PHONY: deploy
deploy:
	gcloud app deploy

.PHONY: attach_log
attach_log:
	gcloud app logs tail -s default