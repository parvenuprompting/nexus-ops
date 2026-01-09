package sys

type ForgeSettings struct {
	Format      string `json:"format"`      // "png", "jpg"
	AspectRatio string `json:"aspectRatio"` // "original", "1:1", "16:9", "9:16", "3:2", "2:3"
}
