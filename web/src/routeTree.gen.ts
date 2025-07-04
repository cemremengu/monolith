/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file was automatically generated by TanStack Router.
// You should NOT make any changes in this file as it will be overwritten.
// Additionally, you should also exclude this file from your linter and/or formatter to prevent it from being checked or modified.

import { Route as rootRouteImport } from './routes/__root'
import { Route as LoginRouteImport } from './routes/login'
import { Route as AuthenticatedRouteImport } from './routes/_authenticated'
import { Route as IndexRouteImport } from './routes/index'
import { Route as AuthenticatedUsersRouteImport } from './routes/_authenticated/users'
import { Route as AuthenticatedSettingsRouteImport } from './routes/_authenticated/settings'
import { Route as AuthenticatedProfileRouteImport } from './routes/_authenticated/profile'
import { Route as AuthenticatedDashboardRouteImport } from './routes/_authenticated/dashboard'

const LoginRoute = LoginRouteImport.update({
  id: '/login',
  path: '/login',
  getParentRoute: () => rootRouteImport,
} as any)
const AuthenticatedRoute = AuthenticatedRouteImport.update({
  id: '/_authenticated',
  getParentRoute: () => rootRouteImport,
} as any)
const IndexRoute = IndexRouteImport.update({
  id: '/',
  path: '/',
  getParentRoute: () => rootRouteImport,
} as any)
const AuthenticatedUsersRoute = AuthenticatedUsersRouteImport.update({
  id: '/users',
  path: '/users',
  getParentRoute: () => AuthenticatedRoute,
} as any)
const AuthenticatedSettingsRoute = AuthenticatedSettingsRouteImport.update({
  id: '/settings',
  path: '/settings',
  getParentRoute: () => AuthenticatedRoute,
} as any)
const AuthenticatedProfileRoute = AuthenticatedProfileRouteImport.update({
  id: '/profile',
  path: '/profile',
  getParentRoute: () => AuthenticatedRoute,
} as any)
const AuthenticatedDashboardRoute = AuthenticatedDashboardRouteImport.update({
  id: '/dashboard',
  path: '/dashboard',
  getParentRoute: () => AuthenticatedRoute,
} as any)

export interface FileRoutesByFullPath {
  '/': typeof IndexRoute
  '/login': typeof LoginRoute
  '/dashboard': typeof AuthenticatedDashboardRoute
  '/profile': typeof AuthenticatedProfileRoute
  '/settings': typeof AuthenticatedSettingsRoute
  '/users': typeof AuthenticatedUsersRoute
}
export interface FileRoutesByTo {
  '/': typeof IndexRoute
  '/login': typeof LoginRoute
  '/dashboard': typeof AuthenticatedDashboardRoute
  '/profile': typeof AuthenticatedProfileRoute
  '/settings': typeof AuthenticatedSettingsRoute
  '/users': typeof AuthenticatedUsersRoute
}
export interface FileRoutesById {
  __root__: typeof rootRouteImport
  '/': typeof IndexRoute
  '/_authenticated': typeof AuthenticatedRouteWithChildren
  '/login': typeof LoginRoute
  '/_authenticated/dashboard': typeof AuthenticatedDashboardRoute
  '/_authenticated/profile': typeof AuthenticatedProfileRoute
  '/_authenticated/settings': typeof AuthenticatedSettingsRoute
  '/_authenticated/users': typeof AuthenticatedUsersRoute
}
export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath
  fullPaths: '/' | '/login' | '/dashboard' | '/profile' | '/settings' | '/users'
  fileRoutesByTo: FileRoutesByTo
  to: '/' | '/login' | '/dashboard' | '/profile' | '/settings' | '/users'
  id:
    | '__root__'
    | '/'
    | '/_authenticated'
    | '/login'
    | '/_authenticated/dashboard'
    | '/_authenticated/profile'
    | '/_authenticated/settings'
    | '/_authenticated/users'
  fileRoutesById: FileRoutesById
}
export interface RootRouteChildren {
  IndexRoute: typeof IndexRoute
  AuthenticatedRoute: typeof AuthenticatedRouteWithChildren
  LoginRoute: typeof LoginRoute
}

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/login': {
      id: '/login'
      path: '/login'
      fullPath: '/login'
      preLoaderRoute: typeof LoginRouteImport
      parentRoute: typeof rootRouteImport
    }
    '/_authenticated': {
      id: '/_authenticated'
      path: ''
      fullPath: ''
      preLoaderRoute: typeof AuthenticatedRouteImport
      parentRoute: typeof rootRouteImport
    }
    '/': {
      id: '/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof IndexRouteImport
      parentRoute: typeof rootRouteImport
    }
    '/_authenticated/users': {
      id: '/_authenticated/users'
      path: '/users'
      fullPath: '/users'
      preLoaderRoute: typeof AuthenticatedUsersRouteImport
      parentRoute: typeof AuthenticatedRoute
    }
    '/_authenticated/settings': {
      id: '/_authenticated/settings'
      path: '/settings'
      fullPath: '/settings'
      preLoaderRoute: typeof AuthenticatedSettingsRouteImport
      parentRoute: typeof AuthenticatedRoute
    }
    '/_authenticated/profile': {
      id: '/_authenticated/profile'
      path: '/profile'
      fullPath: '/profile'
      preLoaderRoute: typeof AuthenticatedProfileRouteImport
      parentRoute: typeof AuthenticatedRoute
    }
    '/_authenticated/dashboard': {
      id: '/_authenticated/dashboard'
      path: '/dashboard'
      fullPath: '/dashboard'
      preLoaderRoute: typeof AuthenticatedDashboardRouteImport
      parentRoute: typeof AuthenticatedRoute
    }
  }
}

interface AuthenticatedRouteChildren {
  AuthenticatedDashboardRoute: typeof AuthenticatedDashboardRoute
  AuthenticatedProfileRoute: typeof AuthenticatedProfileRoute
  AuthenticatedSettingsRoute: typeof AuthenticatedSettingsRoute
  AuthenticatedUsersRoute: typeof AuthenticatedUsersRoute
}

const AuthenticatedRouteChildren: AuthenticatedRouteChildren = {
  AuthenticatedDashboardRoute: AuthenticatedDashboardRoute,
  AuthenticatedProfileRoute: AuthenticatedProfileRoute,
  AuthenticatedSettingsRoute: AuthenticatedSettingsRoute,
  AuthenticatedUsersRoute: AuthenticatedUsersRoute,
}

const AuthenticatedRouteWithChildren = AuthenticatedRoute._addFileChildren(
  AuthenticatedRouteChildren,
)

const rootRouteChildren: RootRouteChildren = {
  IndexRoute: IndexRoute,
  AuthenticatedRoute: AuthenticatedRouteWithChildren,
  LoginRoute: LoginRoute,
}
export const routeTree = rootRouteImport
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>()
