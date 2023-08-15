package scan

type Image struct {
	ImageName                string            `json:"image_name"`
	ImageNameWithProjectName string            `json:"image_name_in_file"`
	ImageTag                 string            `json:"image_tag"`
	ImageNameInHarbor        map[string]string `json:"image_name_in_harbor"`
}

func NewImage(imageName string, imageNameWithProjectName string, imageTag string) *Image {
	return &Image{ImageName: imageName, ImageNameWithProjectName: imageNameWithProjectName, ImageTag: imageTag}
}
