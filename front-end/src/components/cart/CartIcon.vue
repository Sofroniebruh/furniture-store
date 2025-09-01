<script setup lang="ts">
import { ShoppingCart } from 'lucide-vue-next'
import { useCart } from '@/composables/useCart'

defineEmits<{
  click: []
}>()

const { cartItemsCount, isLoading } = useCart()
</script>

<template>
  <button
    @click="$emit('click')"
    class="relative p-2 hover:bg-gray-100 rounded-full transition-colors"
    :class="{ 'animate-pulse': isLoading }"
  >
    <ShoppingCart class="w-6 h-6" />
    <span
      v-if="cartItemsCount > 0"
      class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center"
    >
      {{ cartItemsCount > 99 ? '99+' : cartItemsCount }}
    </span>
  </button>
</template>