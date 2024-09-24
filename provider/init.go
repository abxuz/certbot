package provider

func init() {
	RegisterFactory("aliyun", FactoryFunc(NewAliyunProvider))
	RegisterFactory("cloudflare", FactoryFunc(NewCloudflareProvider))
}
