import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslation } from "react-i18next";
import { useProfile } from "./api/queries";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { User, Mail, Calendar, Shield, Clock } from "lucide-react";
import { toast } from "sonner";
import { useEffect } from "react";

const profileSchema = z.object({
  name: z
    .string()
    .min(1, "Name is required")
    .max(100, "Name must be less than 100 characters"),
  email: z.string().email("Please enter a valid email address"),
  username: z
    .string()
    .min(3, "Username must be at least 3 characters")
    .max(50, "Username must be less than 50 characters"),
});

type ProfileFormData = z.infer<typeof profileSchema>;

export function ProfilePage() {
  const { data: user } = useProfile();
  const { t } = useTranslation();

  const form = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: user?.name || "",
      email: user?.email || "",
      username: user?.username || "",
    },
  });

  // Reset form when user data changes (after fresh fetch)
  useEffect(() => {
    if (user) {
      form.reset({
        name: user.name || "",
        email: user.email || "",
        username: user.username || "",
      });
    }
  }, [user, form]);

  const onSubmit = (data: ProfileFormData) => {
    try {
      // TODO: Update profile information API call
      console.log("Profile update:", data);
      toast.success(t("profile.messages.updateSuccess"));
    } catch (error) {
      console.error("Failed to update profile:", error);
      toast.error(t("profile.messages.updateError"));
    }
  };

  const getInitials = (name: string, username: string) => {
    if (name && name.length > 0) {
      return name
        .split(" ")
        .map((n) => n[0])
        .join("")
        .toUpperCase()
        .slice(0, 2);
    }
    return username.slice(0, 2).toUpperCase();
  };

  return (
    <div className="p-6">
      <div className="mx-auto max-w-4xl space-y-6">
        <div>
          <h1 className="text-3xl font-bold">{t("profile.title")}</h1>
          <p className="text-muted-foreground">{t("profile.subtitle")}</p>
        </div>

        {/* Profile Header */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center space-x-4">
              <Avatar className="h-20 w-20">
                <AvatarImage
                  src={user.avatar}
                  alt={user.name || user.username}
                />
                <AvatarFallback className="text-lg">
                  {getInitials(user.name || "", user.username)}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-1">
                <h2 className="text-2xl font-semibold">
                  {user.name || user.username}
                </h2>
                <p className="text-muted-foreground">@{user.username}</p>
                <div className="text-muted-foreground flex items-center space-x-4 text-sm">
                  <div className="flex items-center space-x-1">
                    <Calendar className="h-3 w-3" />
                    <span>
                      {t("profile.joined")}{" "}
                      {new Date(user.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                  {user.isAdmin && (
                    <div className="flex items-center space-x-1">
                      <Shield className="h-3 w-3" />
                      <span>Administrator</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Profile Information */}
        <Card>
          <CardHeader>
            <CardTitle>{t("profile.information.title")}</CardTitle>
            <CardDescription>
              {t("profile.information.subtitle")}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-6"
              >
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          {t("profile.information.fullName")}
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder={t(
                              "profile.information.fullNamePlaceholder",
                            )}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="username"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          {t("profile.information.username")}
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder={t(
                              "profile.information.usernamePlaceholder",
                            )}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("profile.information.email")}</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder={t(
                            "profile.information.emailPlaceholder",
                          )}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="flex space-x-2">
                  <Button type="submit">
                    {t("profile.actions.saveChanges")}
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => form.reset()}
                  >
                    {t("profile.actions.cancel")}
                  </Button>
                </div>
              </form>
            </Form>
          </CardContent>
        </Card>

        {/* Account Status */}
        <Card>
          <CardHeader>
            <CardTitle>{t("profile.accountStatus.title")}</CardTitle>
            <CardDescription>
              {t("profile.accountStatus.subtitle")}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <User className="text-muted-foreground h-4 w-4" />
                  <span className="text-sm font-medium">
                    {t("profile.accountStatus.accountType")}
                  </span>
                </div>
                <p className="text-muted-foreground pl-6 text-sm">
                  {user.isAdmin
                    ? t("profile.accountStatus.administrator")
                    : t("profile.accountStatus.regularUser")}
                </p>
              </div>
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Mail className="text-muted-foreground h-4 w-4" />
                  <span className="text-sm font-medium">
                    {t("profile.accountStatus.status")}
                  </span>
                </div>
                <p className="text-muted-foreground pl-6 text-sm">
                  {user.isDisabled
                    ? t("profile.accountStatus.disabled")
                    : t("profile.accountStatus.active")}
                </p>
              </div>
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Calendar className="text-muted-foreground h-4 w-4" />
                  <span className="text-sm font-medium">
                    {t("profile.accountStatus.lastUpdated")}
                  </span>
                </div>
                <p className="text-muted-foreground pl-6 text-sm">
                  {new Date(user.updatedAt).toLocaleDateString()}
                </p>
              </div>
              {user.lastSeenAt && (
                <div className="space-y-2">
                  <div className="flex items-center space-x-2">
                    <Clock className="text-muted-foreground h-4 w-4" />
                    <span className="text-sm font-medium">
                      {t("profile.accountStatus.lastSeen")}
                    </span>
                  </div>
                  <p className="text-muted-foreground pl-6 text-sm">
                    {new Date(user.lastSeenAt).toLocaleDateString()}
                  </p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
