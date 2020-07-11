# Action

## Models


### `itemInfoGetReq` 淘宝客商品详情查询（简版）

Name|Type|JSON|Doc
:---|:---|:---|:--
`NumIDs`|`string`|`num_iids`| 商品ID串，用,分割，最大40个
`Platform`|`int`|`platform,omitempty`|链接形式：1：PC，2：无线，默认：１
`Ip	`|`string`|`ip,omitempty`|ip地址，影响邮费获取，如果不传或者传入不准确，邮费无法精准提供

### `itemInfoGetResp` 淘宝客商品详情查询（简版）

Name|Type|JSON|Doc
:---|:---|:---|:--
`TbkItemInfoGetResponse`|`tbkItemInfoGetResponse`|`tbk_item_info_get_response`| 

### `tbkItemInfoGetResponse` 淘宝客商品详情查询（简版）

Name|Type|JSON|Doc
:---|:---|:---|:--
`Results`|`itemInfoGetResults`|`results`| 

### `itemInfoGetResults` 淘宝客商品详情查询（简版）

Name|Type|JSON|Doc
:---|:---|:---|:--
`NTbkItem`|`[]nTbkItem`|`n_tbk_item`|  淘宝客商品

### `smallImages` 淘宝客商品详情查询（简版）

Name|Type|JSON|Doc
:---|:---|:---|:--
`String`|`[]string`|`string`|  淘宝客商品

### `nTbkItem` 淘宝客商品详情查询（简版）

Name|Type|Json|Doc
:---|:---|:---|:--
`CatName                    `|`string      `|`cat_name`|	女装	一级类目名称
`NumIid                     `|`int         `|`num_iid`|123	商品ID
`Title                      `|`string      `|`title`|连衣裙	商品标题
`PictURL                    `|`string      `|`pict_url`|http://gi4.md.alicdn.com/bao/uploaded/i4/xxx.jpg	商品主图
`SmallImages                `|`smallImages `|`small_images`|http://gi4.md.alicdn.com/bao/uploaded/i4/xxx.jpg	商品小图列表
`ReservePrice               `|`string      `|`reserve_price`|102.00	商品一口价格
`ZkFinalPrice               `|`string      `|`zk_final_price`|	88.00	折扣价（元） 若属于预售商品，付定金时间内，折扣价=预售价
`UserType                   `|`int         `|`user_type`|1	卖家类型，0表示集市，1表示商城
`Provcity                   `|`string      `|`provcity`|杭州	商品所在地
`ItemURL                    `|`string      `|`item_url`|http://detail.m.tmall.com/item.htm?id=xxx	商品链接
`SellerID                   `|`int         `|`seller_id`|123	卖家id
`Volume                     `|`int         `|`volume`|1	30天销量
`Nick                       `|`string      `|`nick`|xx旗舰店	店铺名称
`CatLeafName                `|`string      `|`cat_leaf_name`|情趣内衣	叶子类目名称
`IsPrepay                   `|`bool        `|`is_prepay`|true	是否加入消费者保障
`ShopDsr                    `|`int         `|`shop_dsr`|23	店铺dsr 评分
`Ratesum                    `|`int         `|`ratesum`|13	卖家等级
`IRfdRate                   `|`bool        `|`i_rfd_rate`|true	退款率是否低于行业均值
`HGoodRate                  `|`bool        `|`h_good_rate`|true	好评率是否高于行业均值
`HPayRate30                 `|`bool        `|`h_pay_rate30`|true	成交转化是否高于行业均值
`FreeShipment               `|`bool        `|`free_shipment`|true	是否包邮
`MaterialLibType            `|`string      `|`material_lib_type`|1	商品库类型，支持多库类型输出，以英文逗号分隔“,”分隔，1:营销商品主推库，2. 内容商品库，如果值为空则不属于1，2这两种商品类型
`PresaleDiscountFeeText     `|`string      `|`presale_discount_fee_text`|	付定金立减20元	预售商品-商品优惠信息
`PresaleTailEndTime         `|`int64       `|`presale_tail_end_time`|1937297392332	预售商品-付定金结束时间（毫秒）
`PresaleTailStartTime       `|`int64       `|`presale_tail_start_time`|1937297392332	预售商品-付尾款开始时间（毫秒）
`PresaleEndTime             `|`int64       `|`presale_end_time`|	1937297392332	预售商品-付定金结束时间（毫秒）
`PresaleStartTime           `|`int64       `|`presale_start_time`|1937297392332	预售商品-付定金开始时间（毫秒）
`PresaleDeposit             `|`string      `|`presale_deposit`|100	预售商品-定金（元）
`JuPlayEndTime              `|`int64       `|`ju_play_end_time`|1937297392332	聚划算满减 -结束时间（毫秒）    
`JuPlayStartTime            `|`int64       `|`ju_play_start_time`|	1937297392332	聚划算满减 -开始时间（毫秒）
`PlayInfo                   `|`string      `|`play_info`|玩法	1聚划算满减：满N件减X元，满N件X折，满N件X元） 2天猫限时抢：前N分钟每件X元，前N分钟满N件每件X元，前N件每件X元）
`TmallPlayActivityEndTime   `|`int64       `|`tmall_play_activity_end_time`|1937297392332	天猫限时抢可售 -结束时间（毫秒）
`TmallPlayActivityStartTime `|`int64       `|`tmall_play_activity_start_time`|1937297392332	天猫限时抢可售 -开始时间（毫秒）
`JuOnlineStartTime          `|`string      `|`ju_online_start_time`|	1581868800000	聚划算信息-聚淘开始时间（毫秒）
`JuOnlineEndTime            `|`string      `|`ju_online_end_time`|1582300799000	聚划算信息-聚淘结束时间（毫秒）
`JuPreShowStartTime         `|`string      `|`ju_pre_show_start_time`|1581868800000	聚划算信息-商品预热开始时间（毫秒）
`JuPreShowEndTime           `|`string      `|`ju_pre_show_end_time`|1582300799000	聚划算信息-商品预热结束时间（毫秒）
`SalePrice                  `|`string      `|`sale_price`|168	活动价
`KuadianPromotionInfo       `|`string      `|`kuadian_promotion_i`|["每100减20","每200减50"]	跨店满减信息

### `couponGetReq` 阿里妈妈推广券信息查询。传入商品ID+券ID，或者传入me参数，均可查询券信息。

Name|Type|JSON|Doc
:---|:---|:---|:--
`Me`|`string`|`me,omitempty`| nfr%2BYTo2k1PX18gaNN%2BIPkIG2PadNYbBnwEsv6mRavWieOoOE3L9OdmbDSSyHbGxBAXjHpLKvZbL1320ML%2BCF5FRtW7N7yJ056Lgym4X01A%3D	带券ID与商品ID的加密串
`ItemID`|`int`|`item_id,omitempty`|123	商品ID
`ActivityID	`|`string`|`activity_id,omitempty`|sdfwe3eefsdf	券ID

### `couponGetResp` 淘宝客-公用-阿里妈妈推广券详情查询 

Name|Type|JSON|Doc
:---|:---|:---|:--
`TbkCouponGetResponse`|`tbkCouponGetResponse`|`json:"tbk_coupon_get_response`| 

### `tbkCouponGetResponse` 淘宝客-公用-阿里妈妈推广券详情查询 

Name|Type|JSON|Doc
:---|:---|:---|:--
`Data`|`couponGetData`|`json:"data`| 

### `couponGetData` 淘宝客-公用-阿里妈妈推广券详情查询 

Name|Type|JSON|Doc
:---|:---|:---|:--
`CouponStartFee    `|`string `|`coupon_start_fee`|	29.00	优惠券门槛金额
`CouponRemainCount `|`int    `|`coupon_remain_count`|26996	优惠券剩余量
`CouponTotalCount  `|`int    `|`coupon_total_count`|	30000	优惠券总量
`CouponEndTime     `|`string `|`coupon_end_time`|2017-08-17	优惠券结束时间
`CouponStartTime   `|`string `|`coupon_start_time`|	2017-08-15	优惠券开始时间
`CouponAmount      `|`string `|`coupon_amount`|10.00	优惠券金额
`CouponSrcScene    `|`int    `|`coupon_src_scene`|	1	券类型，1 表示全网公开券，4 表示妈妈渠道券
`CouponType        `|`int    `|`coupon_type`|	0	券属性，0表示店铺券，1表示单品券
`CouponActivityID  `|`string `|`coupon_activity_id`|	xsdss	券ID

### `JuTqgGetReq` 获取淘抢购的数据，淘客商品转淘客链接，非淘客商品输出普通链接

Name|Type|JSON|Doc
:---|:---|:---|:--
`AdzoneID`|`int64`|`adzone_id`|123	推广位id（推广位申请方式：http://club.alimama.com/read.php?spm=0.0.0.0.npQdST&tid=6306396&ds=1&page=1&toread=1）
`Fields`|`string`|`fields`|click_url,pic_url,reserve_price,zk_final_price,total_amount,sold_num,title,category_name,start_time,end_time	需返回的字段列表
`StartTime`|`string`|`start_time`|2016-08-09 09:00:00	最早开团时间
`EndTime`|`string`|`end_time`|2016-08-09 16:00:00	最晚开团时间

`PageNO`|`int`|`page_no,omitempty`|1	第几页，默认1，1~100
`PageSize	`|`string`|`page_size,omitempty`|40	页大小，默认40，1~40

### `JuTqgGetResp`  淘抢购api 

Name|Type|JSON|Doc
:---|:---|:---|:--
`TbkJuTqgGetResponse`|`tbkJuTqgGetResponse`|`json:"tbk_ju_tqg_get_response`| 

### `tbkJuTqgGetResponse`  淘抢购api 

Name|Type|JSON|Doc
:---|:---|:---|:--
`Results`|`[]tbkJuTqgGetResults`|`json:"results`| 
`TotalResults`|`int`|`json:"total_results`| 	20	返回的结果数

### `tbkJuTqgGetResults`  淘抢购api 

Name|Type|JSON|Doc
:---|:---|:---|:--
`Title        `|`string `|`title`|连衣裙	商品标题
`TotalAmount  `|`int    `|`total_amount`|100	总库存
`ClickURL     `|`string `|`click_url`|http://s.click.taobao.com/t?e=x	商品链接（是淘客商品返回淘客链接，非淘客商品返回普通h5链接）
`CategoryName `|`string `|`category_name`|潮流女装	类目名称
`ZkFinalPrice `|`string `|`zk_final_price`|50.00	淘抢购活动价
`EndTime      `|`string `|`end_time`|	2016-08-09 13:00:00	结束时间
`SoldNum      `|`int    `|`sold_num`|50	已抢购数量
`StartTime    `|`string `|`start_time`|	2016-08-09 12:00:00	开团时间
`ReservePrice `|`string `|`reserve_price`|	100.00	商品原价
`PicURL       `|`string `|`pic_url`|http: //img4.tbcdn.cn/tfscom/i4/189490253156622336/TB2bZuSsVXXXXcNXXXXXXXXXXXX_!!0-juitemmedia.jpg	商品主图
`NumIid       `|`int    `|`num_iid`|123	商品ID

### `tbkTpwdCreateReq` 提供淘客生成淘口令接口，淘客提交口令内容、logo、url等参数，生成淘口令关键key如：￥SADadW￥，后续进行文案包装组装用于传播

Name|Type|JSON|Doc
:---|:---|:---|:--
`UserID`|`string`|`user_id,omitemptyv`|123	生成口令的淘宝用户ID
`Text`|`string`|`text`|长度大于5个字符	口令弹框内容
`URL`|`string`|`url	`|https://uland.taobao.com/	口令跳转目标页
`Logo`|`string`|`logo,omitempty`|	https://uland.taobao.com/	口令弹框logoURL
`Ext	`|`string`|`ext,omitempty`|	{}	扩展字段JSON格式

### `tbkTpwdCreateResp`  淘宝客-公用-淘口令生成 

Name|Type|JSON|Doc
:---|:---|:---|:--
`TbkTpwdCreateResponse`|`tbkTpwdCreateResponse`|`json:"tbk_tpwd_create_response`| 

### `tbkTpwdCreateResponse`  淘宝客-公用-淘口令生成 

Name|Type|JSON|Doc
:---|:---|:---|:--
`Data`|`tbkTpwdCreateData`|`json:"tbk_tpdatawd_create_response`| 

### `tbkTpwdCreateData`  淘宝客-公用-淘口令生成 

Name|Type|JSON|Doc
:---|:---|:---|:--
`Model`|`string`|`json:"model`| ￥AADPOKFz￥	password

### `itemClickExtractReq` 从长链接或短链接中解析出open_iid

Name|Type|JSON|Doc
:---|:---|:---|:--
`ClickURL	`|`string`|`json:"click_url	`| https://s.click.taobao.com/***	长链接或短链接

### `itemClickExtractResp` 淘宝客-公用-链接解析出商品id 

Name|Type|JSON|Doc
:---|:---|:---|:--
`TbkItemClickExtractResponse	`|`tbkItemClickExtractResponse`|`json:"tbk_item_click_extract_response	`| 

### `tbkItemClickExtractResponse` 淘宝客-公用-链接解析出商品id 

Name|Type|JSON|Doc
:---|:---|:---|:--
`ItemID	`|`string`|`json:"item_id	`| 	123	商品id
`OpenIid	`|`string`|`json:"open_iid	`| xxxxx	商品混淆id
