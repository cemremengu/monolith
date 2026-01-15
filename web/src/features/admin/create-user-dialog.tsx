import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { useCreateUser, useInviteUsers } from "./api/queries";

const createUserSchema = (t: (key: string) => string) =>
  z
    .object({
      username: z
        .string()
        .min(3, t("admin.users.validation.usernameMinLength"))
        .max(50, t("admin.users.validation.usernameMaxLength")),
      name: z.string().min(1, t("admin.users.validation.nameRequired")),
      email: z
        .string()
        .email({ message: t("admin.users.validation.emailInvalid") }),
      role: z.enum(["admin", "user"]),
      password: z
        .string()
        .min(8, t("admin.users.validation.passwordMinLength")),
      confirmPassword: z.string(),
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: t("admin.users.validation.passwordsMustMatch"),
      path: ["confirmPassword"],
    });

const inviteUsersSchema = (t: (key: string) => string) =>
  z.object({
    emails: z
      .string()
      .min(1, t("admin.users.validation.emailsRequired"))
      .refine(
        (value) => {
          const emails = value
            .split(/[\n,]/)
            .map((e) => e.trim())
            .filter(Boolean);
          return emails.every(
            (email) =>
              z.string().email({ message: "" }).safeParse(email).success,
          );
        },
        {
          message: t("admin.users.validation.invalidEmails"),
        },
      ),
    role: z.enum(["admin", "user"]),
  });

type CreateUserFormData = {
  username: string;
  name: string;
  email: string;
  role: "admin" | "user";
  password: string;
  confirmPassword: string;
};

type InviteUsersFormData = {
  emails: string;
  role: "admin" | "user";
};

type CreateUserDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function CreateUserDialog({
  open,
  onOpenChange,
}: CreateUserDialogProps) {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<"invite" | "create">("invite");

  const createUserMutation = useCreateUser();
  const inviteUsersMutation = useInviteUsers();

  const createUserForm = useForm<CreateUserFormData>({
    resolver: zodResolver(createUserSchema(t)),
    defaultValues: {
      username: "",
      name: "",
      email: "",
      role: "user",
      password: "",
      confirmPassword: "",
    },
  });

  const inviteUsersForm = useForm<InviteUsersFormData>({
    resolver: zodResolver(inviteUsersSchema(t)),
    defaultValues: {
      emails: "",
      role: "user",
    },
  });

  const onCreateUser = (data: CreateUserFormData) => {
    createUserMutation.mutate(
      {
        username: data.username,
        name: data.name,
        email: data.email,
        password: data.password,
        isAdmin: data.role === "admin",
      },
      {
        onSuccess: () => {
          toast.success(t("admin.users.messages.userCreated"));
          createUserForm.reset();
          onOpenChange(false);
        },
        onError: () => {
          toast.error(t("admin.users.messages.createUserFailed"));
        },
      },
    );
  };

  const onInviteUsers = (data: InviteUsersFormData) => {
    const emails = data.emails
      .split(/[\n,]/)
      .map((e) => e.trim())
      .filter(Boolean);

    inviteUsersMutation.mutate(
      {
        emails,
        isAdmin: data.role === "admin",
      },
      {
        onSuccess: (response) => {
          if (response.success.length > 0) {
            toast.success(
              t("admin.users.messages.usersInvited", {
                count: response.success.length,
              }),
            );
          }
          if (response.failed.length > 0) {
            response.failed.forEach((failure) => {
              toast.error(`${failure.email}: ${failure.reason}`);
            });
          }
          inviteUsersForm.reset();
          if (response.failed.length === 0) {
            onOpenChange(false);
          }
        },
        onError: () => {
          toast.error(t("admin.users.messages.inviteUsersFailed"));
        },
      },
    );
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      createUserForm.reset();
      inviteUsersForm.reset();
    }
    onOpenChange(newOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[550px]">
        <DialogHeader>
          <DialogTitle>{t("admin.users.dialog.title")}</DialogTitle>
          <DialogDescription>
            {t("admin.users.dialog.description")}
          </DialogDescription>
        </DialogHeader>

        <Tabs
          value={activeTab}
          onValueChange={(v) => setActiveTab(v as "invite" | "create")}
        >
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="invite">
              {t("admin.users.dialog.tabs.invite")}
            </TabsTrigger>
            <TabsTrigger value="create">
              {t("admin.users.dialog.tabs.create")}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="invite" className="space-y-4">
            <Form {...inviteUsersForm}>
              <form
                onSubmit={inviteUsersForm.handleSubmit(onInviteUsers)}
                className="space-y-4"
              >
                <FormField
                  control={inviteUsersForm.control}
                  name="emails"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("admin.users.form.emails")}</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t("admin.users.form.emailsPlaceholder")}
                          className="min-h-[120px]"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={inviteUsersForm.control}
                  name="role"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("admin.users.form.role")}</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue
                              placeholder={t("admin.users.form.selectRole")}
                            />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="user">
                            {t("admin.users.roles.user")}
                          </SelectItem>
                          <SelectItem value="admin">
                            {t("admin.users.roles.admin")}
                          </SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  className="w-full"
                  disabled={inviteUsersMutation.isPending}
                >
                  {inviteUsersMutation.isPending
                    ? t("admin.users.form.inviting")
                    : t("admin.users.form.inviteUsers")}
                </Button>
              </form>
            </Form>
          </TabsContent>

          <TabsContent value="create" className="space-y-4">
            <Form {...createUserForm}>
              <form
                onSubmit={createUserForm.handleSubmit(onCreateUser)}
                className="space-y-4"
              >
                <FormField
                  control={createUserForm.control}
                  name="username"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("admin.users.form.username")}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t(
                            "admin.users.form.usernamePlaceholder",
                          )}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={createUserForm.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("admin.users.form.name")}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("admin.users.form.namePlaceholder")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={createUserForm.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("admin.users.form.email")}</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder={t("admin.users.form.emailPlaceholder")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={createUserForm.control}
                  name="role"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("admin.users.form.role")}</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue
                              placeholder={t("admin.users.form.selectRole")}
                            />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="user">
                            {t("admin.users.roles.user")}
                          </SelectItem>
                          <SelectItem value="admin">
                            {t("admin.users.roles.admin")}
                          </SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={createUserForm.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("admin.users.form.password")}</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder={t(
                            "admin.users.form.passwordPlaceholder",
                          )}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={createUserForm.control}
                  name="confirmPassword"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t("admin.users.form.confirmPassword")}
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder={t(
                            "admin.users.form.confirmPasswordPlaceholder",
                          )}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  className="w-full"
                  disabled={createUserMutation.isPending}
                >
                  {createUserMutation.isPending
                    ? t("admin.users.form.creating")
                    : t("admin.users.form.createUser")}
                </Button>
              </form>
            </Form>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}
