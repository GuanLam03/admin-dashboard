import React from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'
import DefaultLayout from './components/DefaultLayout';
import GuestLayout from './components/GuestLayout';

import DashboardPage from './pages/DashboardPage';
import LoginPage from './pages/LoginPage';
import NotFoundPage from './pages/NotFoundPage';
import ProfilePage from './pages/ProfilePage';
import RolePage from './pages/RolePage';
import PermissionPage from './pages/PermissionPage';
import UserManagementPage from './pages/UserManagementPage';
import UserPermissionPage from './pages/UserPermissionPage';
import DocumentPage from './pages/DocumentPage';
import GoogleDocumentPage from './pages/google_document/GoogleDocumentPage';
import GoogleDocumentAddPage from './pages/google_document/GoogleDocumentAddPage';
import GoogleDocumentEditPage from './pages/google_document/GoogleDocumentEditPage';
import GoogleDocumentViewPage from './pages/google_document/GoogleDocumentViewPage';

import SchedulePage from './pages/schedule/SchedulePage';
import ScheduleAddPage from './pages/schedule/ScheduleAddPage';
import ScheduleEditPage from './pages/schedule/ScheduleEditPage';
import RegisterPage from './pages/RegisterPage';
import SettingsPage from './pages/settings/SettingsPage';
import TwoFactorSettingsPage from './pages/settings/TwoFactorEnablePage';
import TwoFactorDisablePage from './pages/settings/TwoFactorDisablePage';
import EmailsLayout from './pages/email/EmailsLayout';
import TechnicalEmailsPage from './pages/email/TechnicalEmailsPage';
import SupportEmailsPage from './pages/email/SupportEmailsPage';
import MessageViewPage from './pages/email/MessageViewPage';
import EmailsSettingsPage from './pages/email/email_settings/EmailsSettingsPage';
import AdsTrackingPage from './pages/ads_tracking/ads_campaign/AdsCampaignPage';
import AdsCampaignPage from './pages/ads_tracking/ads_campaign/AdsCampaignPage';
import AdsCampaignAddPage from './pages/ads_tracking/ads_campaign/AdsCampaignAddPage';
import AdsTrackingLayout from './pages/ads_tracking/AdsTrackingLayout';
import AdsCampaignEditPage from './pages/ads_tracking/ads_campaign/AdsCampaignEditPage';
import AdsReportPage from './pages/ads_tracking/ads_report/AdsReportPage';
import AdsCampaignReportPage from './pages/ads_tracking/ads_campaign/AdsCampaignReportPage';
import AdsLogPage from './pages/ads_tracking/ads_log/AdsLogPage';








const router = createBrowserRouter([
    {
        path: '/',
        element: <DefaultLayout />,
        children:[
            
            {
                path: '/dashboard',
                element: <DashboardPage />
                
            },
            { path: '/profile', element: <ProfilePage /> },
            
            { path: '/user-management', element: <UserManagementPage /> },
            { path: '/user-management/edit/:id', element: <UserPermissionPage /> },

            { path: '/roles', element: <RolePage /> },
            { path: '/roles/edit/:id', element: <PermissionPage /> },
            { path: '/documents', element: <DocumentPage /> },

            { path: '/google-documents', element: <GoogleDocumentPage /> },
            { path: '/google-documents/add', element: <GoogleDocumentAddPage /> },
            { path: '/google-documents/edit/:id', element: <GoogleDocumentEditPage /> },
            { path: '/google-documents/view/:id', element: <GoogleDocumentViewPage /> },

            { path: '/schedules', element: <SchedulePage /> },
            { path: '/schedules/create', element: <ScheduleAddPage /> },
            { path: '/schedules/edit/:id', element: <ScheduleEditPage /> },

            {
                path: "/emails",
                element: <EmailsLayout />,
                children: [
                    {
                        path: "technical", // /emails/technical
                        element: <TechnicalEmailsPage />,
                    },
                    {
                        path: "support", // /emails/support
                        element: <SupportEmailsPage />,
                    },
                    {
                        path: ":folder/:id", // e.g. /emails/technical/123
                        element: <MessageViewPage />, // detail page
                    },
                    {
                        path: "settings", // e.g. /emails/technical/123
                        element: <EmailsSettingsPage />, // detail page
                    },

                    
                ],
            },
            {
                path: "/ads-tracking",
                element: <AdsTrackingLayout />,
                children: [
                    { path: 'campaign', element: <AdsCampaignPage />, },
                    { path: 'campaign/add', element: <AdsCampaignAddPage /> },
                    { path: 'campaign/edit/:id', element: <AdsCampaignEditPage /> },
                    { path: 'campaign/report/:id', element: <AdsCampaignReportPage /> },


                    { path: 'log', element: <AdsLogPage /> },


                    { path: 'report', element: <AdsReportPage /> },



                    
                ],
            },
            


   

            { path: '/settings', element: <SettingsPage /> },
            { path: '/settings/twofactor', element: <TwoFactorSettingsPage /> },
            { path: '/settings/twofactor/disable', element: <TwoFactorDisablePage /> },



           






            // { path: '/settings', element: <AdminCalendar /> },
            
        ]
    },
    {
        path:'/',
        element: <GuestLayout />,
        children:[
            {
                path: '/login',
                element: <LoginPage />
            },
            {
                path: '/register',
                element: <RegisterPage />
            },
            
        ]
    },
    {
        path: '*',
        element: <NotFoundPage />
    },
])



export default router;
