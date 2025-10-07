<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# Cách sử dụng Jira API để lấy thông tin Development và Deployment của Issue

Dựa trên nghiên cứu sâu từ nhiều nguồn, đây là hướng dẫn toàn diện về cách sử dụng Jira API để lấy thông tin development và deployment của các issue trong Jira.

## **Tổng quan về Development và Deployment trong Jira**

Jira Software Cloud cung cấp hai loại thông tin quan trọng về development:[^1_1][^1_2]

- **Development Information**: Bao gồm commits, branches, pull requests liên quan đến issue
- **Deployment Information**: Bao gồm thông tin về việc deploy code, pipeline, environment


## **Phương pháp chính để truy xuất dữ liệu Development**

### **1. Sử dụng Internal dev-status API (Không được hỗ trợ chính thức)**

Đây là API nội bộ mà Jira UI sử dụng để hiển thị thông tin trong Development panel:[^1_3][^1_4]

#### **Lấy tổng quan Development của Issue:**

```bash
GET /rest/dev-status/latest/issue/summary?issueId=<ISSUE_ID>
```


#### **Lấy chi tiết Development của Issue:**

```bash
GET /rest/dev-status/latest/issue/detail?issueId=<ISSUE_ID>&applicationType=<APPLICATION_TYPE>&dataType=<DATA_TYPE>
```

**Các giá trị dataType hỗ trợ:**

- `repository` - Lấy tất cả thông tin bao gồm commits[^1_5]
- `pullrequest` - Lấy thông tin pull requests
- `branch` - Lấy thông tin branches
- `build` - Lấy thông tin builds

**Ví dụ cụ thể:**

```bash
# Lấy tất cả commits và repository info
curl -u "email@domain.com:API_TOKEN" \
  "https://your-domain.atlassian.net/rest/dev-status/latest/issue/detail?issueId=10195&applicationType=GitHub&dataType=repository"
```

**Lưu ý quan trọng về dev-status API:**

- Chỉ hỗ trợ Basic Authentication với email và API token[^1_4][^1_6]
- Không hỗ trợ OAuth authentication[^1_6]
- Là API nội bộ, không được hỗ trợ chính thức và có thể thay đổi[^1_4]
- Chỉ hoạt động trong môi trường web browser session với Jira[^1_6]


### **2. Sử dụng Development Information REST API (Chính thức)**

Đây là API chính thức để tương tác với development information:[^1_2]

#### **Submit Development Information:**

```bash
POST /rest/devinfo/0.10/bulk
```


#### **Retrieve Repository Information:**

```bash
GET /rest/devinfo/0.10/repository/{repositoryId}
```

**Ví dụ JSON payload để submit development info:**

```json
{
  "repositories": [
    {
      "id": "c6c7c750-cee2-48e2-b920-d7706dfd11f9",
      "name": "atlassian-connect-jira-example",
      "url": "https://bitbucket.org/atlassianlabs/atlassian-connect-jira-example",
      "commits": [
        {
          "id": "a7727ee6350c33cdf90826dc21abaa26a5704370",
          "issueKeys": ["ABC-123"],
          "message": "ABC-123 Update link in documentation",
          "author": {
            "name": "Jane Doe",
            "email": "jane_doe@atlassian.com"
          },
          "authorTimestamp": "2016-10-31T23:27:25+00:00",
          "files": [
            {
              "path": "/README.md",
              "changeType": "MODIFIED",
              "linesAdded": 0,
              "linesRemoved": 1
            }
          ]
        }
      ]
    }
  ]
}
```


## **Phương pháp để truy xuất dữ liệu Deployment**

### **1. Deployment REST API (Chính thức)**

Jira Software Cloud cung cấp Deployment API chính thức:[^1_7]

#### **Submit Deployment Data:**

```bash
POST /rest/deployments/0.1/bulk
```


#### **Get Deployment by Key:**

```bash
GET /rest/deployments/0.1/pipelines/{pipelineId}/environments/{environmentId}/deployments/{deploymentSequenceNumber}
```

**Ví dụ JSON payload để submit deployment info:**

```json
{
  "deployments": [
    {
      "deploymentSequenceNumber": 100,
      "updateSequenceNumber": 1,
      "associations": [
        {
          "associationType": "issueIdOrKeys",
          "values": ["ABC-123", "ABC-456"]
        }
      ],
      "displayName": "Deployment number 16 of Data Depot",
      "url": "http://mydeployer.com/project1/1111-222-333/prod-east",
      "description": "The bits are being transferred",
      "lastUpdated": "2018-01-20T23:27:25.000Z",
      "state": "successful",
      "pipeline": {
        "id": "e9c906a7-451f-4fa6-ae1a-c389e2e2d87c",
        "displayName": "Data Depot Deployment",
        "url": "http://mydeployer.com/project1"
      },
      "environment": {
        "id": "8ec94d72-a4fc-4ac0-b31d-c5a595f373ba",
        "displayName": "US East",
        "type": "production"
      }
    }
  ]
}
```


### **2. Lấy Deployment Information từ Issue**

Hiện tại không có API chính thức để lấy deployment information trực tiếp từ issue. Deployment API chỉ cho phép:[^1_8]

- Submit deployment data
- Get deployment by specific key (pipelineId + environmentId + deploymentSequenceNumber)
- Delete deployment data


## **Authentication Methods**

### **1. Basic Authentication với API Token (Khuyến nghị cho Cloud)**

```bash
# Tạo base64 encoding của email:api_token
echo -n "email@domain.com:your_api_token" | base64

# Sử dụng trong request
curl -H "Authorization: Basic <encoded_string>" \
  "https://your-domain.atlassian.net/rest/api/3/issue/ABC-123"
```


### **2. OAuth 2.0 (Cho integrations phức tạp)**

**Bước 1: Lấy Cloud ID**

```bash
export CLOUD_ID=$(curl "https://your-domain.atlassian.net/_edge/tenant_info" | jq -r '.cloudId')
```

**Bước 2: Lấy Access Token**

```bash
export ACCESS_TOKEN=$(curl --request POST 'https://api.atlassian.com/oauth/token' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "audience": "api.atlassian.com",
    "grant_type": "client_credentials",
    "client_id": "YOUR_CLIENT_ID",
    "client_secret": "YOUR_CLIENT_SECRET"
  }' | jq -r '.access_token')
```

**Bước 3: Sử dụng Token**

```bash
curl --request POST "https://api.atlassian.com/jira/deployments/0.1/cloud/${CLOUD_ID}/bulk" \
  --header "Authorization: Bearer ${ACCESS_TOKEN}" \
  --header 'Content-Type: application/json' \
  --data-raw '{"deployments": [...]}'
```


## **Ví dụ thực tế với Node.js**

### **Lấy Development Information:**

```javascript
const axios = require('axios');
const btoa = require('btoa');

async function getIssueDevelopmentInfo(issueId, domain, email, apiToken) {
  const auth = btoa(`${email}:${apiToken}`);
  
  try {
    // Lấy tổng quan
    const summaryResponse = await axios.get(
      `https://${domain}.atlassian.net/rest/dev-status/latest/issue/summary?issueId=${issueId}`,
      {
        headers: {
          'Authorization': `Basic ${auth}`,
          'Accept': 'application/json'
        }
      }
    );
    
    console.log('Development Summary:', summaryResponse.data);
    
    // Lấy chi tiết commits
    const detailResponse = await axios.get(
      `https://${domain}.atlassian.net/rest/dev-status/latest/issue/detail?issueId=${issueId}&applicationType=GitHub&dataType=repository`,
      {
        headers: {
          'Authorization': `Basic ${auth}`,
          'Accept': 'application/json'
        }
      }
    );
    
    console.log('Development Details:', detailResponse.data);
    
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
}
```


### **Submit Deployment Information:**

```javascript
async function submitDeploymentInfo(cloudId, accessToken, deploymentData) {
  try {
    const response = await axios.post(
      `https://api.atlassian.com/jira/deployments/0.1/cloud/${cloudId}/bulk`,
      {
        deployments: [deploymentData]
      },
      {
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json'
        }
      }
    );
    
    console.log('Deployment submitted:', response.data);
  } catch (error) {
    console.error('Error submitting deployment:', error.response?.data || error.message);
  }
}
```


## **Automation với Jenkins**

Đây là ví dụ Jenkins pipeline script để tự động submit deployment info:[^1_9]

```groovy
pipeline {
  stages {
    stage('Deploy') {
      steps {
        script {
          // Deploy application logic here
          
          // Submit deployment info to Jira
          def deploymentPayload = [
            deployments: [[
              deploymentSequenceNumber: env.BUILD_NUMBER,
              updateSequenceNumber: 1,
              associations: [[
                associationType: "issueIdOrKeys",
                values: ["${JIRA_ISSUE_KEY}"]
              ]],
              displayName: "${JOB_NAME}",
              url: "${BUILD_URL}",
              description: "Automated deployment from Jenkins",
              lastUpdated: new Date().format("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'"),
              state: "successful",
              pipeline: [
                id: "${JOB_NAME}",
                displayName: "${JOB_NAME}",
                url: "${BUILD_URL}"
              ],
              environment: [
                id: "${ENVIRONMENT_ID}",
                displayName: "${ENVIRONMENT_NAME}",
                type: "production"
              ]
            ]]
          ]
          
          sh """
            curl --request POST "https://api.atlassian.com/jira/deployments/0.1/cloud/${CLOUD_ID}/bulk" \
              --header "Authorization: Bearer ${ACCESS_TOKEN}" \
              --header 'Content-Type: application/json' \
              --data '${groovy.json.JsonOutput.toJson(deploymentPayload)}'
          """
        }
      }
    }
  }
}
```


## **Hạn chế và lưu ý quan trọng**

### **Development Information:**

- Internal dev-status API không được hỗ trợ chính thức và có thể thay đổi[^1_4]
- Chỉ hỗ trợ Basic Authentication
- Không thể sử dụng với OAuth[^1_6]


### **Deployment Information:**

- Không có API để lấy deployment info trực tiếp từ issue[^1_8]
- Chỉ có thể submit và retrieve bằng specific keys
- Yêu cầu OAuth 2.0 cho external integrations[^1_10][^1_11]


### **Permissions:**

- Development Tool APIs yêu cầu `jiraDevelopmentTool` module[^1_2]
- Deployment APIs yêu cầu `jiraDeploymentInfoProvider` module[^1_7]
- Scopes cần thiết: `read:dev-info:jira`, `write:dev-info:jira`, `read:deployment-info:jira`, `write:deployment-info:jira`


## **Best Practices**

1. **Sử dụng API chính thức**: Ưu tiên các API được hỗ trợ chính thức thay vì internal APIs[^1_2][^1_7]
2. **Authentication**: Sử dụng API tokens cho Cloud, OAuth cho integrations phức tạp[^1_12]
3. **Rate Limiting**: Implement proper rate limiting và retry logic
4. **Error Handling**: Xử lý errors và validation responses properly
5. **Monitoring**: Monitor API usage và performance
6. **Documentation**: Theo dõi Atlassian documentation cho updates[^1_7][^1_2]

Hướng dẫn này cung cấp cách tiếp cận toàn diện để làm việc với Jira Development và Deployment APIs, bao gồm cả official và workaround methods để truy xuất thông tin cần thiết.
<span style="display:none">[^1_13][^1_14][^1_15][^1_16][^1_17][^1_18][^1_19][^1_20][^1_21][^1_22][^1_23][^1_24][^1_25][^1_26][^1_27][^1_28][^1_29][^1_30][^1_31][^1_32][^1_33][^1_34][^1_35][^1_36][^1_37][^1_38][^1_39][^1_40][^1_41][^1_42][^1_43][^1_44][^1_45][^1_46][^1_47][^1_48][^1_49][^1_50][^1_51][^1_52][^1_53][^1_54][^1_55][^1_56][^1_57][^1_58][^1_59][^1_60][^1_61][^1_62][^1_63][^1_64][^1_65][^1_66][^1_67][^1_68][^1_69][^1_70][^1_71][^1_72][^1_73][^1_74][^1_75][^1_76][^1_77][^1_78][^1_79][^1_80][^1_81]</span>

<div align="center">⁂</div>

[^1_1]: https://docs.gitlab.com/integration/jira/development_panel/

[^1_2]: https://developer.atlassian.com/cloud/jira/software/rest/api-group-development-information/

[^1_3]: https://stackoverflow.com/questions/67800562/rest-api-access-to-information-in-development-panel

[^1_4]: https://jira.atlassian.com/browse/JSWCLOUD-16901

[^1_5]: https://stackoverflow.com/questions/33952456/get-commits-info-of-a-jira-issue-using-rest-api

[^1_6]: https://community.atlassian.com/forums/Jira-questions/Unable-to-access-rest-dev-status-1-0-issue-detail-endpoint-with/qaq-p/3088727

[^1_7]: https://developer.atlassian.com/cloud/jira/software/rest/api-group-deployments/

[^1_8]: https://community.atlassian.com/forums/Jira-questions/How-to-get-Deployments-info-from-Jira-REST-API/qaq-p/2237953

[^1_9]: https://dev.to/pranavlonsane/streamlining-jenkins-jira-integration-automating-deployment-data-submission-4do7

[^1_10]: https://community.developer.atlassian.com/t/using-deployment-api-rest-deployments-0-1-bulk/85975

[^1_11]: https://community.atlassian.com/forums/Jira-Cloud-Admins-discussions/Submitting-Deployment-data-through-API-via-Postman/td-p/1503624

[^1_12]: https://developer.atlassian.com/cloud/jira/software/basic-auth-for-rest-apis/

[^1_13]: https://www.apwide.com/how-to-use-jira-software-for-deployment-management/

[^1_14]: https://developer.atlassian.com/server/jira/platform/jira-rest-api-examples/

[^1_15]: https://support.atlassian.com/jira/kb/access-git-repo-commit-data-via-jira-rest-api-in-development-panel-data/

[^1_16]: https://community.atlassian.com/forums/Jira-questions/Get-Development-Releases-fields-using-API/qaq-p/2689104

[^1_17]: https://jira.atlassian.com/browse/JSWSERVER-15768

[^1_18]: https://gitlab.com/gitlab-org/gitlab/-/merge_requests/18329

[^1_19]: https://developer.atlassian.com/cloud/jira/platform/rest/v2/

[^1_20]: https://www.postman.com/postman/atlassian-jira-api/folder/otfaoew/deployments

[^1_21]: https://github.com/apache/incubator-devlake/issues/4304

[^1_22]: https://developer.atlassian.com/server/framework/atlassian-sdk/rest-api-development/

[^1_23]: https://developer.atlassian.com/cloud/jira/software/modules/deployment/

[^1_24]: https://help.moveworkforward.com/JIGIT/get-started-guide-for-the-jira-end-user

[^1_25]: https://www.postman.com/postman/atlassian-jira-api/request/wdv6sbf/store-development-information

[^1_26]: https://developer.atlassian.com/cloud/jira/software/rest/

[^1_27]: https://deviniti.com/support/addon/server/testflo-813/latest/rest-api/

[^1_28]: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-status/

[^1_29]: https://github.com/marketplace/actions/jira-upload-deployment-info

[^1_30]: https://stackoverflow.com/questions/25667911/how-can-i-find-the-status-of-a-jira-issue-via-the-rest-api

[^1_31]: https://developer.atlassian.com/server/jira/platform/rest/v11000/

[^1_32]: https://help.gitkraken.com/git-integration-for-jira-data-center/rest-api-gij-self-managed/

[^1_33]: https://www.youtube.com/watch?v=yRglBW7YnjA

[^1_34]: https://www.postman.com/postman/atlassian-jira-api/request/u7kp2q8/submit-deployment-data

[^1_35]: https://www.merge.dev/blog/how-to-get-and-create-issues-with-the-jira-api-with-code-snippets

[^1_36]: https://developer.atlassian.com/server/jira/platform/jira-rest-api-tutorials-6291593/

[^1_37]: https://help.gitkraken.com/git-integration-for-jira-data-center/branches-api-gij-self-managed/

[^1_38]: https://developer.atlassian.com/server/jira/platform/jira-rest-api-example-query-issues-6291606/

[^1_39]: https://community.atlassian.com/forums/Jira-questions/With-curl-how-can-I-determine-what-Jira-I-m-using/qaq-p/2541892

[^1_40]: https://docs.arandasoft.com/aic/en/pages/integracion_asms_jira/Configuracion/getStatus_jira.html

[^1_41]: https://stackoverflow.com/questions/64112726/curl-command-to-fetch-data-of-jira-ticket-using-its-ticket-id

[^1_42]: https://help.gitkraken.com/git-integration-for-jira-data-center/get-commits-gij-self-managed/

[^1_43]: https://developer.atlassian.com/server/jira/platform/rest/v10000/intro/

[^1_44]: https://stackoverflow.com/questions/31052721/creating-jira-issue-using-curl-from-command-line

[^1_45]: https://stackoverflow.com/questions/24065109/getting-jira-issues-branches-from-rest-api

[^1_46]: https://www.getknit.dev/blog/deep-dive-developer-guide-to-building-a-jira-api-integration

[^1_47]: https://developer.atlassian.com/server/jira/platform/oauth/

[^1_48]: https://hevodata.com/learn/jira-api/

[^1_49]: https://developer.atlassian.com/server/jira/platform/jira-rest-api-example-oauth-authentication-6291692/

[^1_50]: https://www.miniorange.com/atlassian/rest-api-authentication-using-authorization-grant-from-oauth-provider/

[^1_51]: https://confluence.atlassian.com/spaces/ADMINJIRASERVER/pages/1115659070/Jira+OAuth+2.0+provider+API

[^1_52]: https://community.developer.atlassian.com/t/permissions-to-get-issue-development-information-commits-pull-requests/5911

[^1_53]: https://developer.atlassian.com/cloud/confluence/oauth-2-3lo-apps/

[^1_54]: https://agiletechnicalexcellence.com/2024/04/07/jira-api-intro-authentication.html

[^1_55]: https://www.youtube.com/watch?v=2ahOpLVcYqQ

[^1_56]: https://stackoverflow.com/questions/50723733/jira-dev-status-api-not-returning-all-files-commited

[^1_57]: https://www.miniorange.com/atlassian/rest-api-authentication-in-atlassian-using-custom-provider

[^1_58]: https://www.scribd.com/document/822055207/JIRA-REST-API-Example-Create-Issue-7897248

[^1_59]: https://www.youtube.com/watch?v=gQn0FLe9x-0

[^1_60]: https://community.atlassian.com/forums/Jira-articles/DevOps-way-to-track-CI-CD-Deployments-and-Agile-Release/ba-p/2165788

[^1_61]: https://www.youtube.com/watch?v=E6uSMyU0d_U

[^1_62]: https://developer.atlassian.com/cloud/jira/software/getting-started-open-devops/

[^1_63]: https://developer.atlassian.com/cloud/jira/service-desk/modules/deployment/

[^1_64]: https://www.tempo.io/blog/jira-export

[^1_65]: https://www.postman.com/postman/atlassian-jira-api/request/6dg9jsr/get-a-deployment-by-key

[^1_66]: https://www.linkedin.com/pulse/epic-story-jira-data-extraction-aneesh-menon

[^1_67]: https://support.atlassian.com/jira/kb/troubleshoot-the-development-panel-in-jira-server/

[^1_68]: https://zuplo.com/learning-center/jira-api

[^1_69]: https://community.atlassian.com/forums/Jira-questions/Would-it-be-possible-to-extract-repository-information-for-a-set/qaq-p/1379033

[^1_70]: https://developer.atlassian.com/cloud/jira/software/integrate-jsw-cloud-with-onpremises-tools/

[^1_71]: https://moldstud.com/articles/p-comparing-jira-rest-api-with-other-project-management-apis-a-comprehensive-guide

[^1_72]: https://developer.atlassian.com/server/jira/platform/security-overview/

[^1_73]: https://community.atlassian.com/forums/Jira-questions/How-to-authenticate-to-Jira-REST-API/qaq-p/814987

[^1_74]: https://community.atlassian.com/forums/Jira-questions/When-to-use-OAuth-2-0-and-when-Personal-Access-Token-PAT/qaq-p/1784577

[^1_75]: https://www.miniorange.com/atlassian/atlassian-rest-api-authentication

[^1_76]: https://community.developer.atlassian.com/t/rest-api-authentication-methods/20947

[^1_77]: https://community.developer.atlassian.com/t/creating-jira-issues-oauth-access-token-vs-api-key/75151

[^1_78]: https://developer.atlassian.com/developer-guide/auth/

[^1_79]: https://github.com/sooperset/mcp-atlassian

[^1_80]: https://www.miniorange.com/blog/ways-to-secure-atlassian-rest-api/

[^1_81]: https://forum.uipath.com/t/oauth-2-0-authentication-vs-token-authentication-to-use-jira-application-scope-activity/327807

