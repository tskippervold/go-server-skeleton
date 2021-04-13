#!/bin/sh

gcloud builds submit --tag gcr.io/pay-285612/zo-core

gcloud run deploy zo-core --image gcr.io/pay-285612/zo-core --platform managed --region europe-north1
