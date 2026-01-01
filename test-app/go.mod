module github.com/Sound-X-Team/x-pay/integrations/sdks/golang/test-app

go 1.19

require (
	github.com/Sound-X-Team/x-pay/integrations/sdks/golang v0.0.0
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.5.1
	github.com/shopspring/decimal v1.4.0
)

replace github.com/Sound-X-Team/x-pay/integrations/sdks/golang => ../
