package larksdk

type WikiSpace struct {
	SpaceID     string `json:"space_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SpaceType   string `json:"space_type"`
	Visibility  string `json:"visibility"`
	OpenSharing string `json:"open_sharing"`
}

type WikiNode struct {
	SpaceID         string `json:"space_id"`
	NodeToken       string `json:"node_token"`
	ObjToken        string `json:"obj_token"`
	ObjType         string `json:"obj_type"`
	ParentNodeToken string `json:"parent_node_token"`
	NodeType        string `json:"node_type"`
	OriginNodeToken string `json:"origin_node_token"`
	OriginSpaceID   string `json:"origin_space_id"`
	HasChild        bool   `json:"has_child"`
	Title           string `json:"title"`
	ObjCreateTime   string `json:"obj_create_time"`
	ObjEditTime     string `json:"obj_edit_time"`
	NodeCreateTime  string `json:"node_create_time"`
	Creator         string `json:"creator"`
	Owner           string `json:"owner"`
	NodeCreator     string `json:"node_creator"`
}

type WikiMember struct {
	MemberType string `json:"member_type"`
	MemberID   string `json:"member_id"`
	MemberRole string `json:"member_role"`
	Type       string `json:"type"`
}

type WikiMoveResult struct {
	Node      WikiNode `json:"node"`
	Status    int      `json:"status"`
	StatusMsg string   `json:"status_msg"`
}

type WikiTask struct {
	TaskID      string           `json:"task_id"`
	MoveResults []WikiMoveResult `json:"move_results"`
}

type WikiV1Node struct {
	NodeID     string  `json:"node_id"`
	SpaceID    string  `json:"space_id"`
	ParentID   string  `json:"parent_id"`
	ObjType    int     `json:"obj_type"`
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	Icon       string  `json:"icon"`
	AreaID     string  `json:"area_id"`
	SortID     float64 `json:"sort_id"`
	Domain     string  `json:"domain"`
	ObjToken   string  `json:"obj_token"`
	CreateTime string  `json:"create_time"`
	UpdateTime string  `json:"update_time"`
	DeleteTime string  `json:"delete_time"`
	ChildNum   int     `json:"child_num"`
	Version    int     `json:"version"`
}

type ListWikiSpacesRequest struct {
	PageSize  int
	PageToken string
}

type ListWikiSpacesResult struct {
	Items     []WikiSpace
	PageToken string
	HasMore   bool
}

type GetWikiSpaceRequest struct {
	SpaceID string
	Lang    string
}

type GetWikiNodeRequest struct {
	Token   string
	ObjType string
}

type ListWikiNodesRequest struct {
	SpaceID         string
	ParentNodeToken string
	PageSize        int
	PageToken       string
}

type ListWikiNodesResult struct {
	Items     []WikiNode
	PageToken string
	HasMore   bool
}

type ListWikiSpaceMembersRequest struct {
	SpaceID   string
	PageSize  int
	PageToken string
}

type ListWikiSpaceMembersResult struct {
	Members   []WikiMember
	PageToken string
	HasMore   bool
}

type CreateWikiSpaceMemberRequest struct {
	SpaceID          string
	MemberType       string
	MemberID         string
	MemberRole       string
	Type             string
	NeedNotification bool
}

type DeleteWikiSpaceMemberRequest struct {
	SpaceID    string
	MemberType string
	MemberID   string
	MemberRole string
	Type       string
}

type GetWikiTaskRequest struct {
	TaskID   string
	TaskType string
}

type SearchWikiNodesRequest struct {
	Query     string
	SpaceID   string
	NodeID    string
	PageSize  int
	PageToken string
}

type SearchWikiNodesResult struct {
	Items     []WikiV1Node
	PageToken string
	HasMore   bool
}
