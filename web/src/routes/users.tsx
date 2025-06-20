import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useUsers, useCreateUser, useUpdateUser, useDeleteUser } from '@/lib/queries'
import type { User, CreateUserRequest } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/users')({
  component: Users,
})

function Users() {
  const [formData, setFormData] = useState<CreateUserRequest>({ username: '', name: '', email: '' })
  const [editingUser, setEditingUser] = useState<User | null>(null)

  const { data: users = [], isLoading, error } = useUsers()
  const createUser = useCreateUser()
  const updateUser = useUpdateUser()
  const deleteUser = useDeleteUser()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (editingUser) {
      updateUser.mutate(
        { id: editingUser.id, data: formData },
        {
          onSuccess: () => {
            setFormData({ username: '', name: '', email: '' })
            setEditingUser(null)
          }
        }
      )
    } else {
      createUser.mutate(formData, {
        onSuccess: () => {
          setFormData({ username: '', name: '', email: '' })
        }
      })
    }
  }

  const handleEdit = (user: User) => {
    setEditingUser(user)
    setFormData({ username: user.username, name: user.name || '', email: user.email })
  }

  const handleDelete = (id: string) => {
    deleteUser.mutate(id)
  }

  const handleCancel = () => {
    setEditingUser(null)
    setFormData({ username: '', name: '', email: '' })
  }

  if (isLoading) {
    return <div className="p-6">Loading...</div>
  }

  if (error) {
    return <div className="p-6">Error loading users: {error.message}</div>
  }

  return (
    <div className="p-6">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Users</h1>
        
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>{editingUser ? 'Edit User' : 'Add New User'}</CardTitle>
            <CardDescription>
              {editingUser ? 'Update user information' : 'Create a new user'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Input
                  placeholder="Username"
                  value={formData.username}
                  onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                  required
                />
              </div>
              <div>
                <Input
                  placeholder="Name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                />
              </div>
              <div>
                <Input
                  type="email"
                  placeholder="Email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  required
                />
              </div>
              <div className="flex gap-2">
                <Button 
                  type="submit" 
                  disabled={createUser.isPending || updateUser.isPending}
                >
                  {createUser.isPending || updateUser.isPending 
                    ? 'Saving...' 
                    : editingUser ? 'Update' : 'Create'
                  }
                </Button>
                {editingUser && (
                  <Button type="button" variant="outline" onClick={handleCancel}>
                    Cancel
                  </Button>
                )}
              </div>
            </form>
          </CardContent>
        </Card>

        <div className="grid gap-4">
          {users.map((user) => (
            <Card key={user.id}>
              <CardContent className="p-4">
                <div className="flex justify-between items-center">
                  <div>
                    <h3 className="font-semibold">{user.name || user.username}</h3>
                    <p className="text-sm text-muted-foreground">@{user.username}</p>
                    <p className="text-sm text-muted-foreground">{user.email}</p>
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
                      {deleteUser.isPending ? 'Deleting...' : 'Delete'}
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
              <p className="text-muted-foreground">No users found. Create your first user above.</p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}