<script setup>
import * as z from "zod";
import {useField, useForm} from "vee-validate";
import {toTypedSchema} from "@vee-validate/zod";
import {Button} from '@/components/ui/button'

const schema = z.object({
  email: z.string().min(1, "Email is required").email("Email is not valid").default(""),
  password: z.string().min(1, "Password is required").min(6, "Password should be at least 6 characters").default(""),
})

const {handleSubmit, meta, values} = useForm({
  validationSchema: toTypedSchema(schema),
})

const {
  value: email,
  errorMessage: emailError,
} = useField('email')
const {
  value: password,
  errorMessage: passwordError,
} = useField('password')

const onSubmit = handleSubmit(async values => {
  console.log(values)
})
</script>

<template>
  <div class="flex flex-col gap-8">
    <h1 class="font-semibold text-2xl text-center mt-4">FUMI</h1>
    <form @submit.prevent="onSubmit" class="flex flex-col w-full items-center gap-4">
      <div class="flex flex-col w-full">
        <label for="email">Email</label>
        <input class="border rounded-sm px-2 py-1" v-model="email" id="email" placeholder="example@email.com"/>
        <span class="text-red-500 text-sm">{{ emailError }}</span>
      </div>
      <div class="flex flex-col w-full">
        <label for="email">Password</label>
        <input class="border rounded-sm px-2 py-1" type="password" v-model="password" id="password" placeholder="yourpassword"/>
        <span class="text-red-500 text-sm">{{ passwordError }}</span>
      </div>
      <div class="flex justify-between items-center w-full">
        <Button type="button" class="cursor-pointer" variant="outline">Forgot password</Button>
        <Button type="submit" class="cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]">Sign In</Button>
      </div>
    </form>
  </div>
</template>