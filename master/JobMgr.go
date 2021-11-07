package master

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/api/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"new-crontab/common"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMrg() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   []string{},
		DialTimeout: 5000 * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

// 保存任务到etcd
func (jobMrg *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	// 把任务保存到 /cron/jobs/任务名 -> json
	var (
		jobKey       string
		jobValue     []byte
		putResp      *clientv3.PutResponse
		oldJobObject common.Job
	)

	jobKey = "/cron/jobs/" + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	if putResp, err = jobMrg.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	if putResp.PrevKv == nil {
		return
	}
	if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObject); err != nil {
		err = nil
		return
	}
	oldJob = &oldJobObject
	return
}

func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
	)

	jobKey = "/cron/jobs/" + name

	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	if len(delResp.PrevKvs) == 0 {
		return
	}

	if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
		err = nil
		return
	}

	oldJob = &oldJobObj

	return
}

func (jobMgr *JobMgr) ListJobs() (jobList []*common.Job, err error) {

	var (
		jobKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)

	jobKey = "/cron/jobs/"

	if getResp, err = jobMgr.kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, &job); err != nil {
			err = nil
			return
		}
		jobList = append(jobList, job)
	}

	return
}

func (jobMgr *JobMgr) KillJob(name string) (err error) {

	var (
		killKey        string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseID        clientv3.LeaseID
	)

	killKey = "/cron/kill/" + name

	// 让worker监听到一次put操作，创建一个租约让其稍后自动过期即可
	if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	leaseID = leaseGrantResp.ID

	// 设置killer标记
	if _, err = jobMgr.kv.Put(context.TODO(), killKey, "", clientv3.WithLease(leaseID)); err != nil {
		return
	}

	return
}
