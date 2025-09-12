export const reviews = [
    {
        username: "Bob",
        rating: 5,
        review: "The buying process was seamless, and the furniture arrived on time. The craftsmanship is excellent, and every piece feels sturdy and well-made. I especially love the finish on the dining table—it's both elegant and durable!"
    },
    {
        username: "Josh",
        rating: 5,
        review: "I'm so happy with my purchase! The sofa I ordered is not only stylish but incredibly comfortable. It complements my living room perfectly, and I've already received so many compliments on it. Definetely coming back!"
    },
    {
        username: "Alice",
        rating: 5,
        review: "Great experience overall! The customer service team was helpful in answering my questions, and the delivery was quick. The quality of the materials used in the bed frame is top-notch—I know it's built to last."
    }, {
        username: "Mike",
        rating: 5,
        review: "This is my second time buying from here, and once again, I'm impressed. The furniture is beautiful, and the attention to detail is evident. The chairs I got are both sturdy and lightweight, perfect for everyday use."
    },
]

export const events = [
    {
        name: "Sale"
    },
    {
        name: "None"
    },
    {
        name: "Featured"
    },
    {
        name: "New"
    },
    {
        name: "Best seller"
    }
]

export const productsData = [
    {
        id: "03d419bc-c0e1-4844-a193-6b700d55c514",
        name: "Cool lamp 1",
        description: "Modern stylish lamp perfect for any room",
        amount: 15,
        price: 40.00,
        pictureUrls: ["/images/furniture1.webp"],
        event: "none",
        model: "lamp",
        colors: [
            {
                id: "69cf7cb3-74f1-488a-acb6-a34a3a716b2d",
                name: "White"
            },
            {
                id: "082e1aeb-220b-4b92-ab5c-5a94ee1a94c3",
                name: "Black"
            },
            {
                id: "f3e8d9c2-1a4b-5c6d-7e8f-9a0b1c2d3e4f",
                name: "Gold"
            }
        ]
    },
    {
        id: "b4e7d8c3-2f5a-4b6c-8d9e-1a2b3c4d5e6f",
        name: "Cool lamp 2",
        description: "Elegant desk lamp with adjustable brightness",
        amount: 8,
        price: 40.00,
        pictureUrls: ["/images/furniture2.webp"],
        event: "sale",
        model: "lamp",
        colors: [
            {
                id: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
                name: "Silver"
            },
            {
                id: "b2c3d4e5-f6a7-8901-bcde-f23456789012",
                name: "Bronze"
            },
            {
                id: "c3d4e5f6-a7b8-9012-cdef-345678901234",
                name: "Copper"
            },
            {
                id: "d4e5f6a7-b8c9-0123-def4-56789012345a",
                name: "Chrome"
            }
        ]
    },
    {
        id: "c5f8e9d4-3a6b-5c7d-9e0f-2b3c4d5e6f7a",
        name: "Cool lamp 3",
        description: "Contemporary floor lamp with modern design",
        amount: 12,
        price: 40.00,
        pictureUrls: ["/images/furniture3.webp"],
        event: "none",
        model: "lamp",
        colors: [
            {
                id: "e5f6a7b8-c9d0-1234-ef56-789012345678",
                name: "Matte Black"
            },
            {
                id: "f6a7b8c9-d0e1-2345-f678-90123456789a",
                name: "Wood Finish"
            },
            {
                id: "a7b8c9d0-e1f2-3456-789a-bcdef0123456",
                name: "Navy Blue"
            },
            {
                id: "b8c9d0e1-f2a3-4567-89ab-cdef01234567",
                name: "Forest Green"
            },
            {
                id: "c9d0e1f2-a3b4-5678-9abc-def012345678",
                name: "Burgundy"
            }
        ]
    },
    {
        id: "d6a9f0e5-4b7c-6d8e-0f1a-3c4d5e6f7a8b",
        name: "Cool lamp 4",
        description: "Vintage style table lamp with warm lighting",
        amount: 20,
        price: 40.00,
        pictureUrls: ["/images/furniture4.webp"],
        event: "featured",
        model: "lamp",
        colors: [
            {
                id: "d0e1f2a3-b4c5-6789-abcd-ef0123456789",
                name: "Antique Brass"
            },
            {
                id: "e1f2a3b4-c5d6-789a-bcde-f01234567890",
                name: "Oil Rubbed Bronze"
            },
            {
                id: "f2a3b4c5-d6e7-89ab-cdef-012345678901",
                name: "Pewter"
            },
            {
                id: "a3b4c5d6-e7f8-9abc-def0-123456789012",
                name: "Nickel"
            },
            {
                id: "b4c5d6e7-f8a9-abcd-ef01-234567890123",
                name: "Gunmetal"
            },
            {
                id: "c5d6e7f8-a9b0-bcde-f012-345678901234",
                name: "Rose Gold"
            }
        ]
    },
    {
        id: "e7baef6-5c8d-7e9f-1a2b-4d5e6f7a8b9c",
        name: "Cool lamp 5",
        description: "Industrial pendant lamp for modern spaces",
        amount: 6,
        price: 40.00,
        pictureUrls: ["/images/furniture5.webp"],
        event: "new",
        model: "lamp",
        colors: [
            {
                id: "d6e7f8a9-b0c1-cdef-0123-456789012345",
                name: "Raw Steel"
            },
            {
                id: "e7f8a9b0-c1d2-def0-1234-56789012345a",
                name: "Weathered Iron"
            },
            {
                id: "f8a9b0c1-d2e3-ef01-2345-6789012345ab",
                name: "Rust Finish"
            },
            {
                id: "a9b0c1d2-e3f4-f012-3456-789012345abc",
                name: "Galvanized"
            },
            {
                id: "b0c1d2e3-f4a5-0123-4567-89012345abcd",
                name: "Powder Coated Black"
            },
            {
                id: "c1d2e3f4-a5b6-1234-5678-9012345abcde",
                name: "Industrial Grey"
            },
            {
                id: "d2e3f4a5-b6c7-2345-6789-012345abcdef",
                name: "Copper Patina"
            }
        ]
    },
    {
        id: "f8cbe1a7-6d9e-8f0a-2b3c-5e6f7a8b9c0d",
        name: "Cool lamp 6",
        description: "Minimalist bedside lamp with touch controls",
        amount: 25,
        price: 40.00,
        pictureUrls: ["/images/furniture6.webp"],
        event: "bestseller",
        model: "lamp",
        colors: [
            {
                id: "e3f4a5b6-c7d8-3456-789a-bcdef0123456",
                name: "Pure White"
            },
            {
                id: "f4a5b6c7-d8e9-4567-89ab-cdef01234567",
                name: "Soft Grey"
            },
            {
                id: "a5b6c7d8-e9f0-5678-9abc-def012345678",
                name: "Warm Beige"
            },
            {
                id: "b6c7d8e9-f0a1-6789-abcd-ef0123456789",
                name: "Cream"
            }
        ]
    }
]