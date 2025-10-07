#!/bin/bash

# Load environment variables from .env file
set -a
source "$(dirname "$0")/../../.env"
set +a

JIRA_HOST="${ATLASSIAN_HOST}"
JIRA_EMAIL="${ATLASSIAN_EMAIL}"
JIRA_TOKEN="${ATLASSIAN_TOKEN}"
ISSUE_KEY="SHTP-6050"

echo "========================================="
echo "Step 1: Getting numeric issue ID for ${ISSUE_KEY}"
echo "========================================="

ISSUE_RESPONSE=$(curl -s -u "${JIRA_EMAIL}:${JIRA_TOKEN}" \
  -H "Accept: application/json" \
  "${JIRA_HOST}/rest/api/3/issue/${ISSUE_KEY}?fields=id")

echo "$ISSUE_RESPONSE" | jq '.'

ISSUE_ID=$(echo "$ISSUE_RESPONSE" | jq -r '.id')

if [ -z "$ISSUE_ID" ] || [ "$ISSUE_ID" = "null" ]; then
  echo "ERROR: Could not get issue ID"
  exit 1
fi

echo ""
echo "Issue ID: ${ISSUE_ID}"
echo ""

echo "========================================="
echo "Step 2: Getting development summary"
echo "========================================="

curl -s -u "${JIRA_EMAIL}:${JIRA_TOKEN}" \
  -H "Accept: application/json" \
  "${JIRA_HOST}/rest/dev-status/latest/issue/summary?issueId=${ISSUE_ID}" | jq '.'

echo ""
echo "========================================="
echo "Step 3: Getting development details (repository) - GitLab"
echo "========================================="

curl -s -u "${JIRA_EMAIL}:${JIRA_TOKEN}" \
  -H "Accept: application/json" \
  "${JIRA_HOST}/rest/dev-status/latest/issue/detail?issueId=${ISSUE_ID}&applicationType=GitLab&dataType=repository" | jq '.'

echo ""
echo "========================================="
echo "Step 4: Getting development details (branch) - GitLab"
echo "========================================="

curl -s -u "${JIRA_EMAIL}:${JIRA_TOKEN}" \
  -H "Accept: application/json" \
  "${JIRA_HOST}/rest/dev-status/latest/issue/detail?issueId=${ISSUE_ID}&applicationType=GitLab&dataType=branch" | jq '.'

echo ""
echo "========================================="
echo "Step 5: Getting development details (pullrequest) - GitLab"
echo "========================================="

curl -s -u "${JIRA_EMAIL}:${JIRA_TOKEN}" \
  -H "Accept: application/json" \
  "${JIRA_HOST}/rest/dev-status/latest/issue/detail?issueId=${ISSUE_ID}&applicationType=GitLab&dataType=pullrequest" | jq '.'
