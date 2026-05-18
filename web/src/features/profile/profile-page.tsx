import { zodResolver } from "@hookform/resolvers/zod";
import { User, Mail, Calendar, Shield, Clock } from "lucide-react";
import { useCallback, useEffect } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { z } from "zod";

import { FormInput } from "@/components/form/controlled";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

import { useProfile } from "./api/queries";

const profileSchema = z.object({
  name: z.string().min(1, "Name is required").max(100, "Name must be less than 100 characters"),
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
      toast.success(t("profile.messages.updateSuccess") + data.name);
    } catch (error) {
      console.error("Failed to update profile:", error);
      toast.error(t("profile.messages.updateError"));
    }
  };

  const handleReset = useCallback(() => form.reset(), [form]);

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
    <div className="px-6 py-3">
      <div className="mx-auto max-w-7xl space-y-6">
        <div>
          <h1 className="text-3xl font-bold">{t("profile.title")}</h1>
          <p className="text-muted-foreground">{t("profile.subtitle")}</p>
        </div>

        {/* Profile Header */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center space-x-4">
              <Avatar className="h-20 w-20">
                <AvatarImage src={user.avatar} alt={user.name || user.username} />
                <AvatarFallback className="text-lg">
                  {getInitials(user.name || "", user.username)}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-1">
                <h2 className="text-2xl font-semibold">{user.name || user.username}</h2>
                <p className="text-muted-foreground">@{user.username}</p>
                <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                  <div className="flex items-center space-x-1">
                    <Calendar className="h-3 w-3" />
                    <span>
                      {t("profile.joined")} {new Date(user.createdAt).toLocaleDateString()}
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
            <CardDescription>{t("profile.information.subtitle")}</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <FormInput
                  control={form.control}
                  name="name"
                  label={t("profile.information.fullName")}
                  placeholder={t("profile.information.fullNamePlaceholder")}
                  autoComplete="name"
                />

                <FormInput
                  control={form.control}
                  name="username"
                  label={t("profile.information.username")}
                  placeholder={t("profile.information.usernamePlaceholder")}
                  autoComplete="username"
                />
              </div>

              <FormInput
                control={form.control}
                name="email"
                type="email"
                label={t("profile.information.email")}
                placeholder={t("profile.information.emailPlaceholder")}
                autoComplete="email"
              />

              <div className="flex space-x-2">
                <Button type="submit">{t("profile.actions.saveChanges")}</Button>
                <Button type="button" variant="outline" onClick={handleReset}>
                  {t("profile.actions.cancel")}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        {/* Account Status */}
        <Card>
          <CardHeader>
            <CardTitle>{t("profile.accountStatus.title")}</CardTitle>
            <CardDescription>{t("profile.accountStatus.subtitle")}</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <User className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">
                    {t("profile.accountStatus.accountType")}
                  </span>
                </div>
                <p className="pl-6 text-sm text-muted-foreground">
                  {user.isAdmin
                    ? t("profile.accountStatus.administrator")
                    : t("profile.accountStatus.regularUser")}
                </p>
              </div>
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">{t("profile.accountStatus.status")}</span>
                </div>
                <p className="pl-6 text-sm text-muted-foreground">
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
                <p className="pl-6 text-sm text-muted-foreground">
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
                  <p className="pl-6 text-sm text-muted-foreground">
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
