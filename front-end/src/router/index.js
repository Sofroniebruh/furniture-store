import HeroPage from "@/components/HeroPage.vue";
import {createRouter, createWebHistory} from "vue-router";
import AboutUs from "@/components/AboutUs.vue";
import ContactUs from "@/components/ContactUs.vue";
import Products from "@/components/Products.vue";

const routes = [
    {
        path: '/', name: "Hero page", component: HeroPage,
    },
    {
        path: '/about', name: "About Us page", component: AboutUs,
    },
    {
        path: "/contact", name: "Contact Us page", component: ContactUs,
    },
    {
        path: "/products", name: "Products page", component: Products,
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})