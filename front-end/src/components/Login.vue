<script setup>
import * as z from "zod";
import {useField, useForm} from "vee-validate";
import {toTypedSchema} from "@vee-validate/zod";
import {Button} from '@/components/ui/button'
import {useAuthStore} from "@/stores/useAuth.js";
import {toast} from "vue-sonner";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.js";
import {useRouter} from "vue-router";

const schema = z.object({
  email: z.string().min(1, "Email is required").email("Email is not valid").default(""),
  password: z.string().min(1, "Password is required").min(6, "Password should be at least 6 characters").default(""),
})

const auth = useAuthStore()
const router = useRouter()
const {setOpenDialog, setAllSheetsClosed} = useScreenSheetStore()
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
  try {
    const res = await auth.login(values)

    if (res?.status === 200) {
      toast.success("Log in successfully")
      setOpenDialog(false, {name: "LoginDialog"})

      if (window.innerWidth < 640) {
        setAllSheetsClosed()

        if (res.body?.user?.id) {
          await router.push("/profile/" + res.body.user.id)
        }
      }
    } else if (res?.status === 403 && res.body?.error?.includes("not verified")) {
      toast.error("Please verify your email before logging in")
    }
  } catch (error) {
    console.error(error)
    toast.error("An unexpected error occurred")
  }
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
        <input class="border rounded-sm px-2 py-1" type="password" v-model="password" id="password"
               placeholder="yourpassword"/>
        <span class="text-red-500 text-sm">{{ passwordError }}</span>
        <span class="text-red-500 text-sm">{{ auth.error }}</span>
      </div>
      <div class="flex justify-between items-center w-full">
        <Button type="button" class="cursor-pointer" variant="outline">Forgot password</Button>
        <Button type="submit" class="cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]">Sign In</Button>
      </div>
    </form>
  </div>
</template>