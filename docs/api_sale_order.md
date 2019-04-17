## **1. 分类频道**

### **1.1 Usage**
> 用于在首页展示自定义的分类集合，可以拉取一系列分类下的商品

### **1.2 Method GET /api/category/channel**

### **1.3 Params**
> 无

### **1.4 Example**

#### **1.4.1 Request**
>https://www.geluxiya.com/api/category/channel

### **1.4.2 Result**
```json
{
  "data": [
    {
      "id": 1,
      "name": "热带雨林",
      "iconUrl": "https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/mangguo-icon.png"
    },
    {
      "id": 2,
      "name": "媛媛相抱",
      "iconUrl": "https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/shanzhu-icon.png"
    },
    {
      "id": 3,
      "name": "淡淡的忧伤",
      "iconUrl": "https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/yadan-icon.png"
    },
    {
      "id": 4,
      "name": "美颜组合",
      "iconUrl": "https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/yadan-icon.png"
    },
    {
      "id": 5,
      "name": "营养之王",
      "iconUrl": "https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/liulian-icon.png"
    }
  ],
  "errorMsg": "Ok",
  "errno": 0
}
```


