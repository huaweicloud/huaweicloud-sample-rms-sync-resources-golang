# huaweicloud-sample-rms-sync-resources-golang

[![LICENSE](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/huaweicloud/huaweicloud-sample-rms-sync-resources-golang/blob/master/LICENSE)

huaweicloud-sample-rms-sync-resources-golang is a golang example for you to build CMDB based on your HUAWEI CLOUD resources.

## Requirements

To make the example work, you need a HUAWEI COULD account with resources, then you should configure the resource recorder of Resource Management Service(RMS), specifically the SMN topic which you'll subscribe by different methods.

This example uses MySQL 5.7 to store resources.

## Usage

1. Clone the repository

2. Configure the conf/properties.yml

   ```yaml
   # Domain configuration
   domain:
     domainId: "***"
     ak: "***"
     sk: "***"
   
   # RMS configuration
   rms:
     endpoint: "https://rms.myhuaweicloud.com"
     resourceTypes:
       - vpc.vpcs
       - vpc.securityGroups
       - vpc.securityGroups
       - vpc.publicips
       - ecs.cloudservers
       - evs.volumes
     limit: 100
   
   # SMN configuration
   smn:
     items:
       - regionId: cn-east-3
         endpoint: "https://smn.cn-east-3.myhuaweicloud.com:443"
   
   # Mysql
   mysql:
     username: "root"
     password: "***"
     network: "tcp"
     server: "localhost"
     port: 3306
     database: "rms"
   
   # Sync
   sync:
     spec: "@every 24h"
   ```

3. Run the program

## Eventual consistency

In the example, we use two synchronization mechanism to ensure the eventual consistency of your HUAWEI CLOUD resources in CMDB, one is Full-Sync, the other is Delta-Sync. Also we have introduced two fields to check whether a resource should be update or not.

- **Checksum**: A target to check whether a resource has changed. Once any field of a resource changed, its checksum changed.
- **QueryAt**: A timestamp to query the up-to-date resource from RMS.

**Attention**
>
> - **Full-Sync or Delta-Sync, only when the checksum of resource from CMDB is not equal to that from RMS and queryAt of resource from CMDB is before that from RMS can we update the resource**
>
> - **We just pretend to delete a resource which is marking its state as Deleted and updating the queryAt.**

### Full-Sync

Full-Sync runs once a day to sync all your resources on HUAWEI CLOUD to CMDB. Considering that your resources on HUAWEI CLOUD may be large, they will be synchronized at the granularity of provider, type and region.

Let's take your **ecs.cloudservers** resources in region **cn-north-1** as an example. When it comes to Full-Sync job, the program will perform the following things.

1. List all your **ecs.cloudservers** in **cn-north-1** from Resource Management Service as Set A.
2. Query all your **ecs.cloudservers** in **cn-north-1** from CMDB as Set B.
3. Find the intersection and two subtractions which stands for the resources need to be updated, created and deleted.
4. Perform these actions.

Suppose the number of Set A is N and the number of Set B is M, the time complexity of distinguishing state of all resources in a Full-Sync by provider, type and region is O(max(N, M)).

### Delta-Sync

Delta-Sync receives messages from SMN to sync resources in real time.

Once we received a verified legal message, we'll take different actions based on its notification type. In the sample, we just deal with the **ResourceChanged** notification which has three event type **CREATE** , **UPDATE** and **DELETE**. Then, we perform corresponding operations on CMDB.

You may be concerned that two atypical scenarios, message loss and message disorder, will affect the consistency of resources in your CMDB. 

##### Message loss

Let's take a resource as an example and not consider message disorder for now, we will show you how to ensure consistency when different event types of its **ResourceChanged** notification are lost.

| event type | whether it is the last notification? | how to ensure concurrency?                            |
| ---------- | :----------------------------------: | ----------------------------------------------------- |
| CREATE     |                  Y                   | Next Full-Sync                                        |
| CREATE     |                  N                   | The next related notification comes or Next Full-Sync |
| UPDATE     |                  Y                   | Next Full-Sync                                        |
| UPDATE     |                  N                   | The next related notification comes or Next Full-Sync |
| DELETE     |    Must be the last notification     | Next Full-Sync                                        |

##### Message disorder

As we all know, among many ResourceChanged notifications of a resource, there will only be one **CREATE** notification, at most one **DELETE** notification and many **UPDATE** notifications.

Let's take a resource as an example and not consider message loss for now, we will show you how to ensure consistency when **ResourceChanged** notifications are disordered.

- **We use C, U, D to stand for CREATE, UPDATE and DELETE.**

| original order | actual order | how to ensure concurrency?                                   |
| :------------- | :----------- | :----------------------------------------------------------- |
| C->U           | U->C         | U insert the resource into CMDB. When C comes, the checksum is not equal but the queryAt in CMDB is after that in resource from notification, so we don't update the resource in CMDB. |
| U1->U2         | U2->U1       | U2 update the resource in CMDB, when U1 comes, the checksum is not equal but the queryAt in CMDB is after that in resource from notification, so we don't update the resource in CMDB. |
| U->D           | D->U         | D delete the resource in CMDB, when U comes, the query is not after the queryAt in CMDB, so we do nothing. |
| C->D           | D->C         | D do nothing for there is no related resource in CMDB. When C comes, we insert the resource into CMDB. In this situation, only next Full-Sync can ensure the concurrency of this resource. |

## Warning
>
> The cron scheduler in the sample is a stand-alone machine. If your production environment is distributed, please refer to the distributed cron scheduler such as  RichardKnop/machinery.

