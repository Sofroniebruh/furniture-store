<script setup lang="ts">
import {Product} from "@/lib/types";

interface Props {
  product: Product | null
}

const props = defineProps<Props>();
const emit = defineEmits<{
  close: []
  'delete-product': [productId: string]
}>();

const close = () => {
  emit('close');
}

const handleDelete = (productId: string): void => {
  emit('delete-product', productId);
}

</script>

<template>
  <div>
    <div class="px-4 sm:px-6 py-4 border-b border-gray-200">
      <h3 class="text-lg font-medium text-gray-900">Are you sure?</h3>
    </div>
    <div class="p-4 sm:p-6">
      <p class="text-gray-600 mb-4">
        You are about to delete "{{ props.product?.name }}"? This action cannot be undone.
      </p>
      <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
        <button
            @click="close"
            class="flex-1 px-4 py-3 sm:py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
        >
          Cancel
        </button>
        <button
            @click="handleDelete(props.product!.id)"
            class="flex-1 px-4 py-3 sm:py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 font-medium"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
</template>