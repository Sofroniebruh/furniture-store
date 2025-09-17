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
import {useHistoryStore} from "@/stores/useHistory";

const {
  fetchWishlistData,
  wishlistLoading,
  wishlistError,
  historyLoading,
  historyError,
  currentPage,
  totalPages,
  setPage,
  handlePrevious,
  handleNext,
} = useWishlistOrHistory()

const props = defineProps({
  className: {
    required: false
  }
})

const route = useRoute()
const wishlistStore = useWishlistStore()
const historyStore = useHistoryStore()
const {itemsInWishlist} = storeToRefs(wishlistStore)
const {itemsInHistory} = storeToRefs(historyStore)
const {isInitialized, isLoadingWishlist,} = storeToRefs(wishlistStore)
const {isInitializedHistory, isLoadingHistory} = storeToRefs(historyStore)

const isWishlistLoading = computed(() => {
  return wishlistLoading.value || isLoadingWishlist.value || !isInitialized.value
})
const isHistoryLoading = computed(() => {
  return historyLoading.value || isLoadingHistory.value || !isInitializedHistory.value
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
            :class="cn(isOnProfilePage ?
              `
                 h-[calc(100%-185px)] sm:h-full
                 overflow-y-auto
                 flex flex-col sm:grid
                 grid-cols-2 md:grid-cols-4 lg:grid-cols-5
                 gap-8
                 products-block`
            : 'h-[calc(100%-185px)] overflow-y-auto flex flex-col')">
          <div v-if="isWishlistLoading"
               :class="cn(isOnProfilePage
                 ? 'col-span-full flex justify-center items-center min-h-[400px]'
                 : 'flex-1 flex justify-center items-center')">
            <p class="text-base font-semibold text-gray-600">Loading...</p>
          </div>
          <div v-else-if="wishlistError"
               :class="cn(isOnProfilePage
                 ? 'col-span-full flex justify-center items-center min-h-[400px]'
                 : 'flex-1 flex justify-center items-center')">
            <p class="text-base font-semibold text-gray-600">Error on our side...</p>
          </div>
          <div v-else-if="itemsInWishlist.length === 0"
               :class="cn(isOnProfilePage
                 ? 'col-span-full flex justify-center items-center min-h-[400px]'
                 : 'flex-1 flex justify-center items-center')">
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
          <div v-if="isHistoryLoading"
               :class="cn(isOnProfilePage
                 ? 'col-span-full flex justify-center items-center min-h-[400px]'
                 : 'flex-1 flex justify-center items-center')">
            <p class="text-base font-semibold text-gray-600">Loading...</p>
          </div>
          <div v-else-if="historyError"
               :class="cn(isOnProfilePage
                 ? 'col-span-full flex justify-center items-center min-h-[400px]'
                 : 'flex-1 flex justify-center items-center')">
            <p class="text-base font-semibold text-gray-600">Error on our side...</p>
          </div>
          <div v-else-if="itemsInHistory.length === 0"
               :class="cn(isOnProfilePage
                 ? 'col-span-full flex justify-center items-center min-h-[400px]'
                 : 'flex-1 flex justify-center items-center')">
            <p class="text-base font-semibold text-gray-600">No purchases made</p>
          </div>
          <template v-else>
            <ProductCardProfile
                v-if="!isOnProfilePage"
                v-for="(product, index) in itemsInHistory"
                :key="`history-${product.id}-${index}`"
                :is-history-card="true"
                :product="product"
            />
            <ProductCardProfilePage v-else v-for="(product, index) in itemsInHistory"
                                    :key="`history-${product.id}-${index}`" :is-history-card="true" :product="product"/>
          </template>
        </div>
      </TabsContent>
    </Tabs>
    <Pagination :set-page="setPage"
                :current-page="currentPage"
                :handle-next="handleNext"
                :handle-previous="handlePrevious"
                :total-pages="totalPages"
                v-if="isOnProfilePage"/>
  </div>
</template>