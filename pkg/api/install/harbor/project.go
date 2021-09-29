package harbor

//// ArtifactReservePolicy -> 配置项目保留策略对象
//type ArtifactReservePolicy struct {
//	Algorithm string `json:"algorithm"`
//	Rules     []rule `json:"rules"`
//	Trigger   struct {
//		Kind       string `json:"kind"`
//		References struct {
//		} `json:"references"`
//		Settings struct {
//			Cron string `json:"cron"`
//		} `json:"settings"`
//	} `json:"trigger"`
//	Scope struct {
//		Level string `json:"level"`
//		Ref   int    `json:"ref"`
//	} `json:"scope"`
//}
//
//type rule struct {
//	Disabled       bool           `json:"disabled"`
//	Action         string         `json:"action"`
//	ScopeSelectors ScopeSelectors `json:"scope_selectors"`
//	TagSelectors   []TagSelector  `json:"tag_selectors"`
//	Params         struct {
//		LatestPushedK int8 `json:"latestPushedK"`
//	} `json:"params"`
//	Template string `json:"template"`
//}
//
//type TagSelector struct {
//	Kind       string `json:"kind"`
//	Decoration string `json:"decoration"`
//	Pattern    string `json:"pattern"`
//	Extras     string `json:"extras"`
//}
//
//type ScopeSelectors struct {
//	Repository []Repository `json:"repository"`
//}
//
//type Repository struct {
//	Kind       string `json:"kind"`
//	Decoration string `json:"decoration"`
//	Pattern    string `json:"pattern"`
//}
//
//type ProjectMetadata struct {
//	CreationTime       time.Time `json:"creation_time"`
//	CurrentUserRoleId  int       `json:"current_user_role_id"`
//	CurrentUserRoleIds []int     `json:"current_user_role_ids"`
//	CveAllowlist       struct {
//		CreationTime time.Time     `json:"creation_time"`
//		Id           int           `json:"id,omitempty"`
//		Items        []interface{} `json:"items"`
//		ProjectId    int           `json:"project_id"`
//		UpdateTime   time.Time     `json:"update_time"`
//	} `json:"cve_allowlist"`
//	Metadata struct {
//		Public      string `json:"public"`
//		RetentionId string `json:"retention_id"`
//	} `json:"metadata"`
//	Name       string    `json:"name"`
//	OwnerId    int       `json:"owner_id"`
//	OwnerName  string    `json:"owner_name"`
//	ProjectId  int       `json:"project_id"`
//	UpdateTime time.Time `json:"update_time"`
//}
//
//func (hm harborMgr) listProjects() *harborMgr {
//	log.Println("获取project映射关系...")
//
//	page := 1
//	hm.ProjectsMap = make(map[int]ProjectMetadataSub)
//
//	for {
//		url := fmt.Sprintf("http://%s:%s/api/v2.0/projects?page=%d&page_size=15", hm.Harbor.ResolvAddress, hm.Harbor.HttpPort, page)
//		resp, err := get(url, nil, "admin", hm.Harbor.Password)
//		if err != nil {
//			panic(err)
//		}
//
//		var p []ProjectMetadata
//		b, _ := ioutil.ReadAll(resp.Body)
//		json.Unmarshal(b, &p)
//
//		if len(p) == 0 {
//			return &hm
//		}
//		for _, v := range p {
//			fmt.Println(v.Name, v.ProjectId)
//			hm.ProjectsMap[v.ProjectId] = ProjectMetadataSub{
//				Name:        v.Name,
//				RetentionID: v.Metadata.RetentionId,
//			}
//		}
//
//		page++
//	}
//	return &hm
//}
//
//func (hm harborMgr) setProjectsReserveNumPolicy() *harborMgr {
//	log.Println("配置project制品保留策略...")
//	for k, _ := range hm.ProjectsMap {
//		hm.setProjectReserveNumPolicy(k, hm.Harbor.ResolvAddress)
//	}
//	return &hm
//}
//
//func (hm harborMgr) setProjectReserveNumPolicy(projectID int, host string) {
//	if hm.Harbor.ReserveNum <= 0 {
//		return
//	}
//
//	retentionID := hm.ProjectsMap[projectID].RetentionID
//
//	log.Printf("####### -> 配置项目: %s的制品保留策略，项目id为: %d, Retention id为: %s",
//		hm.ProjectsMap[projectID].Name, projectID, retentionID)
//
//	url := fmt.Sprintf("http://%s:%s/api/v2.0/retentions/%s", host, hm.Harbor.HttpPort, retentionID)
//	log.Printf("请求url为: %s", url)
//
//	tagSelector := TagSelector{
//		Kind:       "doublestar",
//		Decoration: "matches",
//		Pattern:    "**",
//		Extras:     "{\"untagged\":false}",
//	}
//
//	repo := Repository{
//		Kind:       "doublestar",
//		Decoration: "repoMatches",
//		Pattern:    "**",
//	}
//
//	scopeSelectors := ScopeSelectors{Repository: []Repository{repo}}
//
//	r := rule{
//		Disabled:       false,
//		Action:         "retain",
//		ScopeSelectors: scopeSelectors,
//		TagSelectors:   []TagSelector{tagSelector},
//		Params: struct {
//			LatestPushedK int8 `json:"latestPushedK"`
//		}{LatestPushedK: hm.Harbor.ReserveNum},
//		Template: "latestPushedK",
//	}
//
//	var rules []rule
//	rules = append(rules, r)
//
//	reserveNumPolicy := ArtifactReservePolicy{
//		Algorithm: "or",
//		Rules:     rules,
//		Trigger: struct {
//			Kind       string   `json:"kind"`
//			References struct{} `json:"references"`
//			Settings   struct {
//				Cron string `json:"cron"`
//			} `json:"settings"`
//		}{Kind: "Schedule", Settings: triggerSettings{Cron: "0 0 * * * *"}},
//		Scope: struct {
//			Level string `json:"level"`
//			Ref   int    `json:"ref"`
//		}{Level: "project", Ref: projectID},
//	}
//
//	b, err := json.Marshal(reserveNumPolicy)
//
//	log.Printf("请求体为：%s\n", string(b))
//
//	if err != nil {
//		panic(err)
//	}
//
//	resp, err := post(url, bytes.NewBuffer(b), "admin", hm.Harbor.Password)
//	if err != nil {
//		panic(err)
//	}
//
//	log.Printf("响应状态码为: %d\n", resp.StatusCode)
//	result, _ := ioutil.ReadAll(resp.Body)
//	log.Println(string(result))
//}
