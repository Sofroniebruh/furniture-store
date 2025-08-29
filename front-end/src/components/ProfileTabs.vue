<script setup lang="ts">
import {Tabs, TabsContent, TabsList, TabsTrigger,} from '@/components/ui/tabs'
import {useWishlistOrHistory} from "@/composables/useWishlistOrHistory";
import {computed, onMounted} from "vue";
import ProductCardProfile from "@/components/product-related/ProductCardProfile.vue";
import {useWishlistStore} from "@/stores/useWishlist";
import {cn} from "@/lib/utils.js";
import {storeToRefs} from "pinia";
import {useRoute} from "vue-router";
import ProductCardProfilePage from "@/components/product-related/ProductCardProfilePage.vue";
import Pagination from "@/components/Pagination.vue";

const {
  productsHistoryData,
  fetchWishlistData,
  wishlistLoading,
  wishlistError,
  historyLoading,
  historyError
} = useWishlistOrHistory()

const props = defineProps({
  className: {
    required: false
  }
})

const route = useRoute()
const wishlistStore = useWishlistStore()
const {itemsInWishlist} = storeToRefs(wishlistStore)
const {isInitialized, isLoadingWishlist} = useWishlistStore()

const historyItems = productsHistoryData

const isWishlistLoading = computed(() => {
  return wishlistLoading.value || isLoadingWishlist || !isInitialized
})
const isHistoryLoading = computed(() => {
  return historyLoading.value
})
const isOnProfilePage = computed(() => route.path.startsWith('/profile'))

onMounted(async () => {
  try {
    await Promise.all([
      fetchWishlistData("wishlist"),
      fetchWishlistData("history")
    ])
  } catch (error) {
    console.error('Error in component mount:', error)
  }
})
</script>

<template>
  <div :class="cn(isOnProfilePage ? 'flex min-h-[732px]' : 'h-screen flex flex-col', props.className)">
    <Tabs default-value="wishlist" :class="cn(isOnProfilePage ? 'flex-1 flex flex-col justify-center'
    : 'flex-1 w-full flex flex-col justify-center')">
      <div class="flex flex-col items-center">
        <TabsList class="grid w-full max-w-[332px] grid-cols-2 flex-shrink-0">
          <TabsTrigger value="wishlist">
            Wishlist
          </TabsTrigger>
          <TabsTrigger value="history">
            Purchase history
          </TabsTrigger>
        </TabsList>
      </div>

      <TabsContent value="wishlist" class="flex-1 flex flex-col">
        <div
            :class="cn(isOnProfilePage ? 'w-full grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-8 products-block'
            : 'h-[calc(100%-185px)] overflow-y-auto flex flex-col')">
          <div v-if="isWishlistLoading" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Loading...</p>
          </div>
          <div v-else-if="wishlistError" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Error on our side...</p>
          </div>
          <div v-else-if="itemsInWishlist.length === 0" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">No products added</p>
          </div>
          <template v-else>
            <ProductCardProfile
                v-if="!isOnProfilePage"
                v-for="(product, index) in itemsInWishlist"
                :key="`history-${product.id}-${index}`"
                :product="product"
            />
            <ProductCardProfilePage v-else v-for="(product, index) in itemsInWishlist"
                                    :key="`history-${product.id}-${index}`" :product="product"/>
          </template>
        </div>
      </TabsContent>

      <TabsContent value="history" class="w-full flex-1 flex flex-col">
        <div
            :class="cn(isOnProfilePage ? 'w-full grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-8 products-block'
            : 'h-[calc(100%-185px)] overflow-y-auto flex flex-col')">
          <div v-if="isHistoryLoading" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Loading...</p>
          </div>
          <div v-else-if="historyError" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Error on our side...</p>
          </div>
          <div v-else-if="historyItems.length === 0" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">No purchases made</p>
          </div>
          <template v-else>
            <ProductCardProfile
                v-if="!isOnProfilePage"
                v-for="(product, index) in historyItems"
                :key="`history-${product.id}-${index}`"
                :is-history-card="true"
                :product="product"
            />
            <ProductCardProfilePage v-else v-for="(product, index) in historyItems"
                                    :key="`history-${product.id}-${index}`" :is-history-card="true" :product="product"/>
          </template>
        </div>
      </TabsContent>
    </Tabs>
    <Pagination v-if="isOnProfilePage"/>
  </div>
</template>