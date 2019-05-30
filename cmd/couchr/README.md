# couchr

Migrates data to CouchDB

Note: run from the top dir if using `.env` file

```bash
go run cmd/couchr/main.go
```


CPD Activity Doc from current system

```json
{
    "activity": {
      "categoryId": 0,
      "categoryName": "",
      "code": "RACP5",
      "description": "• Reading journals and texts\n• Information searches, e.g. Medline\n• Audio/videotapes\n• Web-based learning\n• Other learning activities",
      "id": 24,
      "name": "Other Learning Activities",
      "unitId": 0,
      "unitName": ""
    },
    "category": {
      "code": "",
      "description": "RACP CPD categories.",
      "id": 10,
      "name": "RACP"
    },
    "createdAt": "0001-01-01T00:00:00Z",
    "credit": 2,
    "creditData": {
      "quantity": 2,
      "quantityFixed": false,
      "unitCode": "",
      "unitCredit": 1,
      "unitDescription": "",
      "unitName": "hours"
    },
    "date": "2018-01-08",
    "dateISO": "0001-01-01T00:00:00Z",
    "description": "review papers on cardiovascular genetics for JAHA and JCDD",
    "evidence": true,
    "id": 6717,
    "memberId": 586,
    "type": {
      "id": {
        "Int64": 32,
        "Valid": true
      },
      "name": "Reading journals and texts"
    },
    "updatedAt": "0001-01-01T00:00:00Z"
}
```


New version:

```json
{
    "type": "cpd",
    "created": "2017-02-01T00:00:00Z",
    "updated": "2017-02-01T00:00:00Z"
    
    "domain": "RACP",
    "category": "Other Learning Activities",
    "activity": "Reading journals and texts",
    "unit": "hours",
    "creditPerUnit": 1,
    
    "date": "2018-01-08",
    "description": "review papers on cardiovascular genetics for JAHA and JCDD",
    "evidence": true,
    "quantity": 2,
    "credit": 2
}
```
