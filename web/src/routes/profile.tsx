import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/lib/auth";
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
import { Separator } from "@/components/ui/separator";
import {
  User,
  Mail,
  Calendar,
  Shield,
  Globe,
  Palette,
  Clock,
} from "lucide-react";
import { ThemeSwitcher } from "@/components/theme-switcher";
import { LanguageSwitcher } from "@/components/language-switcher";

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
  language: z.string().optional(),
  theme: z.string().optional(),
  timezone: z.string().optional(),
});

type ProfileFormData = z.infer<typeof profileSchema>;

export const Route = createFileRoute("/profile")({
  component: Profile,
});

function Profile() {
  const { user } = useAuth();
  const { t } = useTranslation();

  const form = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: user?.name || "",
      email: user?.email || "",
      username: user?.username || "",
      language: user?.language || "en",
      theme: user?.theme || "system",
      timezone: user?.timezone || "UTC",
    },
  });

  const onSubmit = async (data: ProfileFormData) => {
    // TODO: Implement profile update API call
    console.log("Profile update:", data);
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

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-sm text-muted-foreground">
          {t("profile.loading")}
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="max-w-4xl mx-auto space-y-6">
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
                <div className="flex items-center space-x-4 text-sm text-muted-foreground">
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
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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

                <Separator />

                <div className="space-y-4">
                  <h3 className="text-lg font-medium">
                    {t("profile.preferences.title")}
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="space-y-2">
                      <div className="flex items-center space-x-2">
                        <Globe className="h-4 w-4" />
                        <span className="text-sm font-medium">
                          {t("profile.preferences.language")}
                        </span>
                      </div>
                      <LanguageSwitcher />
                    </div>
                    <div className="space-y-2">
                      <div className="flex items-center space-x-2">
                        <Palette className="h-4 w-4" />
                        <span className="text-sm font-medium">
                          {t("profile.preferences.theme")}
                        </span>
                      </div>
                      <ThemeSwitcher />
                    </div>
                    <div className="space-y-2">
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4" />
                        <span>{t("profile.preferences.timezone")}</span>
                      </div>
                      <Input
                        placeholder={t(
                          "profile.preferences.timezonePlaceholder",
                        )}
                      />
                    </div>
                  </div>
                </div>

                <div className="flex space-x-2">
                  <Button type="submit">
                    {t("profile.actions.saveChanges")}
                  </Button>
                  <Button type="button" variant="outline">
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
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <User className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">
                    {t("profile.accountStatus.accountType")}
                  </span>
                </div>
                <p className="text-sm text-muted-foreground pl-6">
                  {user.isAdmin
                    ? t("profile.accountStatus.administrator")
                    : t("profile.accountStatus.regularUser")}
                </p>
              </div>
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">
                    {t("profile.accountStatus.status")}
                  </span>
                </div>
                <p className="text-sm text-muted-foreground pl-6">
                  {user.isDisabled
                    ? t("profile.accountStatus.disabled")
                    : t("profile.accountStatus.active")}
                </p>
              </div>
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Calendar className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">
                    {t("profile.accountStatus.lastUpdated")}
                  </span>
                </div>
                <p className="text-sm text-muted-foreground pl-6">
                  {new Date(user.updatedAt).toLocaleDateString()}
                </p>
              </div>
              {user.lastSeenAt && (
                <div className="space-y-2">
                  <div className="flex items-center space-x-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm font-medium">
                      {t("profile.accountStatus.lastSeen")}
                    </span>
                  </div>
                  <p className="text-sm text-muted-foreground pl-6">
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
