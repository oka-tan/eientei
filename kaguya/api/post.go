package api

//Post represents a 4chan post from a thread, including the OP, as it comes from the API
type Post struct {
	No            int64   `json:"no"`
	Resto         int64   `json:"resto"`
	Time          int32   `json:"time"`
	Sticky        *int8   `json:"sticky"`
	Name          *string `json:"name"`
	Trip          *string `json:"trip"`
	ID            *string `json:"id"`
	Capcode       *string `json:"capcode"`
	Country       *string `json:"country"`
	BoardFlag     *string `json:"board_flag"`
	Sub           *string `json:"sub"`
	Com           *string `json:"com"`
	Tim           *int64  `json:"tim"`
	MD5           *string `json:"md5"`
	Filename      *string `json:"filename"`
	Ext           *string `json:"ext"`
	Fsize         *int64  `json:"fsize"`
	W             *int16  `json:"w"`
	H             *int16  `json:"h"`
	TnW           *int16  `json:"tn_w"`
	TnH           *int16  `json:"tn_h"`
	FileDeleted   *int8   `json:"file_deleted"`
	Spoiler       *int8   `json:"spoiler"`
	CustomSpoiler *int8   `json:"custom_spoiler"`
	Tag           *string `json:"tag"`
	Since4Pass    *int16  `json:"since4pass"`
	Archived      *int8   `json:"archived"`
}
