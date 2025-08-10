<script setup lang="ts">
import {Tabs, TabsContent, TabsList, TabsTrigger,} from '@/components/ui/tabs'
import {useWishlistOrHistory} from "@/composables/useWishlistOrHistory";
import {computed, onMounted} from "vue";
import ProductCardProfile from "@/components/ProductCardProfile.vue";
import {useWishlistStore} from "@/stores/useWishlist";

const {
  productsHistoryData,
  fetchWishlistData,
  wishlistLoading,
  wishlistError,
  historyLoading,
  historyError
} = useWishlistOrHistory()

const {itemsInWishlist, isInitialized, isLoadingWishlist} = useWishlistStore()

const wishlistItems = computed(() => {
  return Array.isArray(itemsInWishlist) ? itemsInWishlist : []
})

const historyItems = computed(() => {
  return Array.isArray(productsHistoryData.value) ? productsHistoryData.value : []
})

const isWishlistLoading = computed(() => {
  return wishlistLoading.value || isLoadingWishlist || !isInitialized
})

const isHistoryLoading = computed(() => {
  return historyLoading.value
})

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
  <div class="h-screen flex flex-col">
    <Tabs default-value="wishlist" class="w-[350px] flex-1 flex flex-col">
      <TabsList class="grid w-full grid-cols-2 flex-shrink-0">
        <TabsTrigger value="wishlist">
          Wishlist
        </TabsTrigger>
        <TabsTrigger value="history">
          Purchase history
        </TabsTrigger>
      </TabsList>

      <TabsContent value="wishlist" class="flex-1 flex flex-col">
        <div class="h-[calc(100%-185px)] overflow-y-auto flex flex-col">
          <div v-if="isWishlistLoading" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Loading...</p>
          </div>
          <div v-else-if="wishlistError" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Error on our side...</p>
          </div>
          <div v-else-if="wishlistItems.length === 0" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">No products added</p>
          </div>
          <div v-else class="w-full">
            <ProductCardProfile
                v-for="(product, index) in wishlistItems"
                :key="`wishlist-${product.id}-${index}`"
                :image-src="product.pictureUrls?.[0]"
                :name="product.name"
                :price="product.price"
                :id="product.id"
                :product="product"
            />
          </div>
        </div>
      </TabsContent>

      <TabsContent value="history" class="flex-1 flex flex-col">
        <div class="h-[calc(100%-185px)] overflow-y-auto flex flex-col">
          <div v-if="isHistoryLoading" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Loading...</p>
          </div>
          <div v-else-if="historyError" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">Error on our side...</p>
          </div>
          <div v-else-if="historyItems.length === 0" class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600">No purchases made</p>
          </div>
          <div v-else class="w-full">
            <ProductCardProfile
                v-for="(product, index) in historyItems"
                :key="`history-${product.id}-${index}`"
                :is-history-card="true"
                :image-src="product.pictureUrls?.[0]"
                :product="product"
                :name="product.name"
                :price="product.price"
                :id="product.id"
            />
          </div>
        </div>
      </TabsContent>
    </Tabs>
  </div>
</template>