export interface Product {
    id: string;
    name: string;
    amount: number;
    price: number;
    description: string;
    pictureUrls: string[];
    event: string;
    model: string;
    colors: Colors[];
}

export interface Colors {
    id: string;
    name: string;
}