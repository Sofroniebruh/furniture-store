export interface Product {
    id: string;
    name: string;
    stock: number;
    price: number;
    description: string;
    pictureUrls: string[];
    event: string;
    model: string;
    colors: Color[];
}

export interface Color {
    id?: string;
    name: string;
}