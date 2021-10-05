package querynessus

import "strings"

type FolderCollection struct {
	Folders []Folder `json:"folders"`
}

func (folderCollection FolderCollection) FolderId(folderName string) (id int, exists bool) {
	for _, folder := range folderCollection.Folders {
		if strings.EqualFold(folder.Name, folderName) {
			return folder.Id, true
		}
	}
	return 0, false
}

type Folder struct {
	UnreadCount int    `json:"unread_count"`
	IsCustom    int    `json:"custom"`
	IsDefault   int    `json:"default_tag"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Id          int    `json:"id"`
}
