import HeroPage from "@/components/HeroPage.vue";
import {createRouter, createWebHistory} from "vue-router";
import AboutUs from "@/components/AboutUs.vue";
import ContactUs from "@/components/ContactUs.vue";
import Products from "@/components/Products.vue";
import ProductPage from "@/components/product-related/ProductPage.vue";
import Profile from "@/components/Profile.vue";
import {useAuthStore} from "@/stores/useAuth.js";
import AuthComponent from "@/components/auth/AuthComponent.vue";
import Dashboard from "@/components/dashboard-related/Dashboard.vue";
import CartPage from "@/components/cart/CartPage.vue";
import CheckoutPage from "@/components/checkout/CheckoutPage.vue";
import CheckoutSuccess from "@/components/checkout/CheckoutSuccess.vue";

const routes = [
    {
        path: '/',
        name: "Hero",
        component: HeroPage,
    },
    {
        path: '/auth',
        name: "Authentication",
        component: AuthComponent,
        meta: { redirectIfAuthenticated: true }
    },
    {
        path: "/profile/:id",
        name: "Profile",
        component: Profile,
        meta: { requiresAuth: true },
    },
    {
        path: "/about",
        name: "About",
        component: AboutUs,
    },
    {
        path: "/contact",
        name: "Contact",
        component: ContactUs,
    },
    {
        path: "/products",
        name: "Products",
        component: Products,
    },
    {
        path: "/dashboard",
        name: "Dashboard",
        component: Dashboard,
        meta: {
            requiresAuth: true,
            roles: ['admin']
        }
    },
    {
        path: "/product/:id",
        name: "Product",
        component: ProductPage,
    },
    {
        path: "/cart",
        name: "Cart",
        component: CartPage,
    },
    {
        path: "/checkout",
        name: "checkout",
        component: CheckoutPage,
        meta: { requiresAuth: true },
    },
    {
        path: "/checkout/success",
        name: "checkout-success",
        component: CheckoutSuccess,
        meta: { requiresAuth: true },
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        redirect: '/'
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    if (!authStore.isAuthenticated && (to.meta.requiresAuth || to.meta.roles)) {
        try {
            await authStore.fetchUser()
        } catch (error) {
            console.error('Error fetching user:', error)
        }
    }

    const isAuthenticated = authStore.isAuthenticated
    const userRoles = authStore.user?.roles || []

    if (to.meta.requiresAuth && !isAuthenticated) {
        return next({
            name: 'Authentication',
            query: { redirect: to.fullPath }
        })
    }

    if (to.meta.redirectIfAuthenticated && isAuthenticated) {
        const redirectPath = to.query.redirect || `/profile/${authStore.user.id}`
        return next(redirectPath)
    }

    if (to.meta.roles && to.meta.roles.length > 0) {
        if (!isAuthenticated) {
            return next({
                name: 'Authentication',
                query: { redirect: to.fullPath }
            })
        }

        const hasRequiredRole = to.meta.roles.some(role => userRoles.includes(role))
        if (!hasRequiredRole) {
            return next({ name: 'Hero' })
        }
    }

    next()
})