<script setup lang="ts">
import {useProducts} from "@/composables/useProducts";
import {useRoute} from "vue-router";
import {computed, onBeforeUnmount, onMounted, ref, watch} from "vue";
import Wrapper from "@/components/Wrapper.vue";
import {Button} from "@/components/ui/button";
import {cn} from "@/lib/utils";
import {CheckCircle, ChevronLeft, ChevronRight, Heart, Shield, Truck} from 'lucide-vue-next'
import {useWishlistStore} from "@/stores/useWishlist";
import {useWishlistOrHistory} from "@/composables/useWishlistOrHistory";

const {fetchProductById, product, loading, error} = useProducts()
const {
  isItemInWishlist,
  toggleWishlist,
  removeFromWishlistStore,
  addToWishlistStore,
  saveToLocalStorage
} = useWishlistStore()
const {addToWishlist, removeFromWishlist} = useWishlistOrHistory()
const route = useRoute()

const productId = computed(() => String(route.params.id))
const activeImageIdx = ref<number>(0)
const quantity = ref<number>(1)

const images = computed(() => product.value?.pictureUrls ?? [])
const activeImage = computed(() => images.value[activeImageIdx.value] ?? "")

const hasPrev = computed(() => activeImageIdx.value > 0)
const hasNext = computed(() => activeImageIdx.value < images.value.length - 1)

const isInStock = computed(() => (product.value?.amount ?? 0) > 0)
const maxQty = computed(() => Math.max(product.value?.amount ?? 0, 0))

const formatCurrency = (value: number) =>
    new Intl.NumberFormat('en-US', {style: 'currency', currency: 'USD'}).format(value)

const formattedPrice = computed(() => formatCurrency(product.value?.price ?? 0))

const decQty = () => {
  quantity.value = Math.max(1, quantity.value - 1)
}
const incQty = () => {
  quantity.value = Math.min(maxQty.value || 1, quantity.value + 1)
}

const handleWishlist = async () => {
  if (!product.value) return
  const isCurrentlyInWishlist = isItemInWishlist(productId.value)

  toggleWishlist(product.value)

  if (!isCurrentlyInWishlist) {
    try {
      const res = await addToWishlist(productId.value)

      console.log("Res", res)

      if (res !== 201) {
        removeFromWishlistStore(productId.value)
      }
    } catch (e) {
      removeFromWishlistStore(productId.value)
      console.error(e)
    }
  } else {
    try {
      const res = await removeFromWishlist(productId.value)

      if (res !== 200) {
        addToWishlistStore(product.value)
      }
    } catch (e) {
      addToWishlistStore(product.value)
      console.error(e)
    }
  }

  saveToLocalStorage()
}
const goPrev = () => {
  if (hasPrev.value) activeImageIdx.value -= 1
}
const goNext = () => {
  if (hasNext.value) activeImageIdx.value += 1
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'ArrowLeft') goPrev()
  if (e.key === 'ArrowRight') goNext()
}

onMounted(async () => {
  await fetchProductById(productId.value)
  window.addEventListener('keydown', handleKeydown)
})

watch(() => route.params.id, async (newId) => {
  if (newId) {
    activeImageIdx.value = 0
    quantity.value = 1
    await fetchProductById(String(newId))
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Wrapper class="pt-6">
    <nav aria-label="Breadcrumb" class="mb-6 text-sm text-gray-500">
      <router-link to="/" class="hover:text-gray-700">Home</router-link>
      <span class="mx-2">/</span>
      <router-link to="/products" class="hover:text-gray-700">Products</router-link>
      <span class="mx-2">/</span>
      <span class="text-gray-700" v-if="!loading && product">{{ product.name }}</span>
      <span class="inline-block h-4 w-32 bg-gray-200 animate-pulse rounded align-middle" v-else/>
    </nav>

    <div class="grid grid-cols-1 md:grid-cols-12 gap-8 items-start justify-center">
      <div class="md:col-span-7 lg:col-span-7 xl:col-span-8">
        <div v-if="loading" class="space-y-4">
          <div class="w-full max-w-2xl aspect-square bg-gray-200 animate-pulse rounded"/>
          <div class="grid grid-cols-4 sm:grid-cols-5 md:grid-cols-6 gap-3 max-w-2xl">
            <div v-for="n in 6" :key="n" class="h-20 bg-gray-200 animate-pulse rounded"></div>
          </div>
        </div>
        <div v-else-if="!product">
          <div class="w-full aspect-square bg-gray-100 rounded flex items-center justify-center text-gray-500">
            No image
          </div>
        </div>
        <div v-else class="space-y-4">
          <div class="relative group">
            <img :src="activeImage" :alt="product.name" class="w-full max-w-2xl aspect-square object-cover rounded"/>
            <button aria-label="Previous image" @click="goPrev" :disabled="!hasPrev"
                    :class="cn('absolute left-2 top-1/2 -translate-y-1/2 bg-white/90 hover:bg-white text-gray-700 disabled:opacity-50 rounded-full w-9 h-9 flex items-center justify-center shadow', hasPrev && 'cursor-pointer')">
              <ChevronLeft class="w-5 h-5"/>
            </button>
            <button aria-label="Next image" @click="goNext" :disabled="!hasNext"
                    :class="cn('absolute right-2 top-1/2 -translate-y-1/2 bg-white/90 hover:bg-white text-gray-700 disabled:opacity-50 rounded-full w-9 h-9 flex items-center justify-center shadow', hasNext && 'cursor-pointer') ">
              <ChevronRight class="w-5 h-5"/>
            </button>
          </div>
          <div class="grid grid-cols-4 sm:grid-cols-5 md:grid-cols-6 gap-3 max-w-2xl">
            <button v-for="(img, idx) in images" :key="idx" @click="activeImageIdx = idx"
                    class="relative rounded overflow-hidden ring-2 focus:outline-none focus:ring-2 focus:ring-[#c9a275]"
                    :aria-pressed="activeImageIdx === idx"
                    :class="activeImageIdx === idx ? 'ring-[#c9a275]' : 'ring-transparent'">
              <img :src="img" :alt="`${product.name} thumbnail ${idx+1}`" class="h-20 w-full object-cover"/>
            </button>
          </div>
        </div>
      </div>

      <div class="space-y-6 md:col-span-5 lg:col-span-5 xl:col-span-4 md:pl-2 lg:pl-6 md:sticky md:top-24">
        <div v-if="loading" class="space-y-4">
          <div class="h-8 bg-gray-200 animate-pulse rounded w-2/3"/>
          <div class="h-6 bg-gray-200 animate-pulse rounded w-1/2"/>
          <div class="h-24 bg-gray-200 animate-pulse rounded w-full"/>
          <div class="h-10 bg-gray-200 animate-pulse rounded w-40"/>
        </div>
        <div v-else-if="error || !product" class="text-center py-12">
          <p class="text-red-600">Failed to load product. Please try again.</p>
        </div>
        <div v-else class="space-y-5">
          <div class="flex items-start justify-between gap-4">
            <h1 class="text-2xl md:text-3xl font-semibold">{{ product.name }}</h1>
            <span v-if="isInStock"
                  class="inline-flex items-center gap-1 text-xs md:text-sm text-green-700 bg-green-50 border border-green-200 px-2 py-1 rounded">
              <CheckCircle class="w-4 h-4"/> In stock
            </span>
            <span v-else
                  class="inline-flex items-center gap-1 text-xs md:text-sm text-red-700 bg-red-50 border border-red-200 px-2 py-1 rounded">
              Out of stock
            </span>
          </div>
          <div class="flex items-center gap-2">
            <p class="text-xl md:text-2xl font-semibold">{{ formattedPrice }}</p>
            <span v-if="product.event" class="text-xs rounded bg-[#c9a275]/10 text-[#c9a275] px-2 py-1">{{
                product.event
              }}</span>
          </div>
          <p class="text-gray-600 leading-relaxed text-base md:text-lg">{{ product.description }}</p>

          <div class="space-y-2">
            <p class="font-medium">Available colors</p>
            <div class="flex flex-wrap gap-2">
              <span v-for="c in product.colors" :key="c.id"
                    class="px-2 py-1 text-sm rounded border border-gray-200 text-gray-700">
                {{ c.name }}
              </span>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <div class="inline-flex items-center border rounded">
              <button @click="decQty" aria-label="Decrease quantity" class="px-3 py-2 text-gray-700 disabled:opacity-50"
                      :disabled="quantity <= 1">-
              </button>
              <input type="number" :max="maxQty" min="1" v-model.number="quantity"
                     class="appearance-none [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-inner-spin-button]:m-0 w-16 text-center border-l border-r py-2"/>
              <button @click="incQty" aria-label="Increase quantity" class="px-3 py-2 text-gray-700"
                      :disabled="!isInStock">+
              </button>
            </div>
            <Button class="cursor-pointer" :disabled="!isInStock">Add to cart</Button>
            <Button @click="handleWishlist" variant="outline"
                    class="cursor-pointer"
                    aria-label="Add to wishlist">
              <Heart :class="cn('w-4 h-4', isItemInWishlist(productId) ? 'fill-red-500 text-red-500' : '')"/>
            </Button>
          </div>

          <div class="grid md:grid-cols-1 sm:grid-cols-3 grid-cols-1 gap-3 text-sm text-gray-700 pt-2">
            <div class="flex items-center gap-2 border rounded p-3">
              <Truck class="w-4 h-4"/>
              Free shipping over $100
            </div>
            <div class="flex items-center gap-2 border rounded p-3">
              <Shield class="w-4 h-4"/>
              2-year warranty
            </div>
            <div class="flex items-center gap-2 border rounded p-3">
              <CheckCircle class="w-4 h-4"/>
              30-day returns
            </div>
          </div>
        </div>
      </div>
    </div>
  </Wrapper>
</template>