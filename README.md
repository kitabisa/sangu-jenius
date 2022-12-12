# Jenius Sangu

## DEPRECATED AS OF 12/12/2022
Moved to https://github.com/kitabisa/sangu-client

## Usage blueprint

1. There is a type named `Client` (`jenius.Client`) that should be instantiated through `NewClient` which hold any possible setting to the library.
2. There is a gateway classes which you will be using depending on whether you used. The gateway type need a Client instance.
3. Any activity (pay request, pay status, etc) is done in the gateway level.

## Example

```go
    jeniusclient := jenius.NewClient()
    jeniusclient.JeniusXChannelId = "YOUR-JENIUS-X-CHANNEL-ID"
    jeniusclient.JeniusApiKey = "YOUR-JENIUS-API-KEY"
    jeniusclient.JeniusApiSecret = "YOUR-JENIUS-API-SECRET"
    jeniusclient.JeniusClientId  = "YOUR-JENIUS-CLIENT-ID"
    jeniusclient.JeniusClientSecret = "YOUR-JENIUS-CLIENT-SECRET"
    jeniusclient.JeniusBaseUrl = "YOUR-JENIUS-BASE-URL"s

    coreGateway := jenius.CoreGateway{
        Client: jeniusclient,
    }

    payReq := &jenius.PayRequestReq{
        ReferenceNo: "190324101010700A",
        Token: "YOUR-TOKEN",
    }

    payReqBody := &jenius.PayRequestReqBody{
        TxnAmount:    "20000",
        Cashtag:      "$cashTag",
        PromoCode:    "PROMOCODE",
        UrlCallback:  "YOUR-JENIUS-URL-CALLBACK",
        PurchaseDesc: "Description",
        CreatedAt: 1554588000, //Transaction Time
    }

    resp, _, _ := coreGateway.PayRequest(payReq, payReqBody)
```

## Credits
[Midtrans Library for Go(lang)]: https://github.com/veritrans/go-midtrans
