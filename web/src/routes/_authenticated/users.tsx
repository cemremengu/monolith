import { createFileRoute, ErrorComponent } from "@tanstack/react-router";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  useCreateUser,
  useUpdateUser,
  useDeleteUser,
  usersQueryOptions,
} from "@/api/users/queries";
import type { User } from "@/api/users/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
import { useSuspenseQuery } from "@tanstack/react-query";
import { Spinner } from "@/components/spinner";

const userSchema = z.object({
  username: z
    .string()
    .min(3, "Username must be at least 3 characters")
    .max(50, "Username must be less than 50 characters"),
  name: z
    .string()
    .min(1, "Name is required")
    .max(100, "Name must be less than 100 characters"),
  email: z.string().email("Please enter a valid email address"),
});

type UserFormData = z.infer<typeof userSchema>;

export const Route = createFileRoute("/_authenticated/users")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(usersQueryOptions({})),
  component: Users,
  pendingComponent: Spinner,
  errorComponent: ErrorComponent,
});

function Users() {
  const { data: users } = useSuspenseQuery(usersQueryOptions({}));

  const [editingUser, setEditingUser] = useState<User | null>(null);

  const form = useForm<UserFormData>({
    resolver: zodResolver(userSchema),
    defaultValues: {
      username: "",
      name: "",
      email: "",
    },
  });

  const createUser = useCreateUser();
  const updateUser = useUpdateUser();
  const deleteUser = useDeleteUser();

  const onSubmit = (data: UserFormData) => {
    if (editingUser) {
      updateUser.mutate(
        { id: editingUser.id, data },
        {
          onSuccess: () => {
            form.reset();
            setEditingUser(null);
          },
        },
      );
    } else {
      createUser.mutate(data, {
        onSuccess: () => {
          form.reset();
        },
      });
    }
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    form.reset({
      username: user.username,
      name: user.name || "",
      email: user.email,
    });
  };

  const handleDelete = (id: string) => {
    deleteUser.mutate(id);
  };

  const handleCancel = () => {
    setEditingUser(null);
    form.reset();
  };

  return (
    <div className="p-6">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Users</h1>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle>{editingUser ? "Edit User" : "Add New User"}</CardTitle>
            <CardDescription>
              {editingUser ? "Update user information" : "Create a new user"}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-4"
              >
                <FormField
                  control={form.control}
                  name="username"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Username</FormLabel>
                      <FormControl>
                        <Input placeholder="Username" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Name</FormLabel>
                      <FormControl>
                        <Input placeholder="Name" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input type="email" placeholder="Email" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div className="flex gap-2">
                  <Button
                    type="submit"
                    disabled={createUser.isPending || updateUser.isPending}
                  >
                    {createUser.isPending || updateUser.isPending
                      ? "Saving..."
                      : editingUser
                        ? "Update"
                        : "Create"}
                  </Button>
                  {editingUser && (
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleCancel}
                    >
                      Cancel
                    </Button>
                  )}
                </div>
              </form>
            </Form>
          </CardContent>
        </Card>

        <div className="grid gap-4">
          {users.map((user) => (
            <Card key={user.id}>
              <CardContent className="p-4">
                <div className="flex justify-between items-center">
                  <div>
                    <h3 className="font-semibold">
                      {user.name || user.username}
                    </h3>
                    <p className="text-sm text-muted-foreground">
                      @{user.username}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {user.email}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      Created: {new Date(user.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleEdit(user)}
                    >
                      Edit
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDelete(user.id)}
                      disabled={deleteUser.isPending}
                    >
                      {deleteUser.isPending ? "Deleting..." : "Delete"}
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {users.length === 0 && (
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-muted-foreground">
                No users found. Create your first user above.
              </p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
