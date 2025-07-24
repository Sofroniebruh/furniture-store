import {useRoute, useRouter} from "vue-router";

export const useParams = () => {
    const route = useRoute()
    const router = useRouter()

    const updatePage = (page: number) => {
        router.push({
            query: {
                ...route.query,
                page: page.toString(),
            }
        })
    }
    const updatePriceFrom = (price: number) => {
        router.push({
            query: {
                ...route.query,
                price_from: price.toString(),
            }
        })
    }
    const updatePriceTo = (price: number) => {
        router.push({
            query: {
                ...route.query,
                price_to: price.toString(),
            }
        })
    }
    const updateEvent = (event: string) => {
        router.push({
            query: {
                ...route.query,
                event: event,
            }
        })
    }
    const updateModel = (models: []) => {
        router.push({
            query: {
                ...route.query,
                model: models,
            }
        })
    }

    return {
        updateEvent,
        updatePage,
        updatePriceFrom,
        updatePriceTo,
        updateModel,
    }
}