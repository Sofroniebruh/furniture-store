<script setup lang="ts">
import {Tabs, TabsContent, TabsList, TabsTrigger,} from '@/components/ui/tabs'
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
        <div class="h-[calc(100%-185px)] overflow-y-auto">
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
          <p>Wishlist</p>
        </div>
      </TabsContent>
      <TabsContent value="history" class="flex-1 flex flex-col">
        <div class="h-[calc(100%-185px)] overflow-y-auto">
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
          <p>History</p>
        </div>
      </TabsContent>
    </Tabs>
  </div>
</template>