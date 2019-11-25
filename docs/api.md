# BlockQuiz API

## System

### Healthy Check

```http request
GET /hc
```

**Response:**

```json5
{
  "data": {
    "duration": "4m25.558410944s"
  }
}
```

## Task

### Create Task

```http request
POST /task
```

**Param:**

```json5
{
  "language": "en", // en or zh
  "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7" // mixin id
}
```

**Response:**

```json5
{
  "data": {
    "id": 26,
    "created_at": 1573466214,
    "updated_at": 1573466214,
    "language": "en",
    "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
    "creator": "lucky coin",
    "course": {
      "id": 1,
      "language": "en",
      "title": "李永乐老师的比特币入门课（1）",
      "summary": "2008年，网络极客中本聪提出了比特币的概念，这是一种全新的电子货币。比特币是一种去中心化的记账系统，人们通过挖矿获得比特币，通过公开记账的方式完成支付。前几年，比特币价格暴涨数百万倍，让许多人一夜暴富。去年，比特币的暴跌又让许多人损失惨重。我准备通过两期节目给大家介绍比特币的基本原理。在这一期节目中，我将介绍比特币和区块链的基本概念，以及矿机在挖矿时到底在做什么。有兴趣的小朋友，点开视频看看吧！"
    },
    "total_question": 3,
    "state": "PENDING",
    "is_blocked": false,
    "block_until": 1573466214
  }
}
```

### Active Task

```http request
POST /task/:id/active
```

**Response:**

```json5
{
  "data": {
    "id": 26,
    "created_at": 1573466214,
    "updated_at": 1573466344,
    "language": "en",
    "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
    "creator": "lucky coin",
    "state": "COURSE",
    "block_until": 1573466214
  }
}
```

### Task Detail

```http request
GET /task/:id
```

**Response:**

```json5
{
  "data": {
    "id": 26,
    "created_at": 1573466214,
    "updated_at": 1573466344,
    "language": "en",
    "title": "2019-11-26",
    "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
    "creator": "lucky coin",
    "state": "PENDING",
    "block_until": 1573466214
  }
}
```

### Cancel Task

```http request
POST /task/:id/cancel
```

**Response:**

```json5
{
  "data": {
    "id": 26,
    "created_at": 1573466214,
    "updated_at": 1573467064,
    "language": "en",
    "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
    "creator": "lucky coin",
    "state": "CANCELLED",
    "block_until": 1573466214
  }
}
```

## Error

```json5
{
  "error": {
    "code": 10001,
    "msg": "invalid parameters",
    "hint": "Object.language: enhh does not validate as in(en|zh): invalid parameters (10001)" // debug mode 才有
  }
}
```

## Task State

| state | description |
|:------:|:------------:|
|   PENDING    |   未完成 |
| FINISH | 已完成 |
