
# inventory

`inventory` has been written as a small utility addressing specific needs for GRNET's teleconferencing service. It is a fairly lightweight API that can be used to add BigBlueButton and Transcoding types of hosts into Ansible's dynamic inventory and to deploy upon them the necessary Ansible roles. 

In general three methods are implemented:

Name | Verb | Body | Description 
---- | ---- | ---- | -----------
/list | GET | - | List all BigBlueButton and Transcoding servers
/add | POST | application/json | Add a new pair of BigBlueButton and Transcoding servers into Ansible's dynamic inventory
/deploy | POST | - | Deploy via Ansible the whole conferencing instrastructure

Examples:
```bash
curl –X GET 'http://{{ host_or_ip }}:{{ port }}/list'
```

```bash
curl -X POST -H "Content-Type: application/json" \
-d '[{ 
        "name": "webconf-bbb5.grnet.gr", 
        "type": "bigbluebutton", 
        "internal_ip": "172.16.0.52", 
        "storage_path": 5 
    }, 
    { 
        "name": "webconf-trans5.grnet.gr",
        "type": "transcoding", 
        "internal_ip": "172.16.0.51",
        "storage_path": 5
    }]' 'http://{{ host_or_ip }}:{{ port }}/list'
```

```bash
curl –X POST 'http://{{ host_or_ip }}:{{ port }}/deploy'
```
