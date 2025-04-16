import React from 'react';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '../radix/tabs';
import { cn } from '../../../utils/cn';

const ShadcnTabs = Tabs;

const ShadcnTabsList = React.forwardRef<
  React.ElementRef<typeof TabsList>,
  React.ComponentPropsWithoutRef<typeof TabsList>
>(({ className, ...props }, ref) => (
  <TabsList
    ref={ref}
    className={cn(
      'inline-flex h-10 items-center justify-center rounded-md bg-gray-100 dark:bg-gray-800 p-1 text-gray-500 dark:text-gray-400',
      className
    )}
    {...props}
  />
));
ShadcnTabsList.displayName = 'ShadcnTabsList';

const ShadcnTabsTrigger = React.forwardRef<
  React.ElementRef<typeof TabsTrigger>,
  React.ComponentPropsWithoutRef<typeof TabsTrigger> & {
    variant?: 'default' | 'outline' | 'emerald';
  }
>(({ className, variant = 'default', ...props }, ref) => {
  const variantStyles = {
    default: 'data-[state=active]:bg-white dark:data-[state=active]:bg-gray-900 data-[state=active]:text-gray-900 dark:data-[state=active]:text-gray-50',
    outline: 'data-[state=active]:bg-transparent data-[state=active]:border-b-2 data-[state=active]:border-gray-900 dark:data-[state=active]:border-gray-50 rounded-none px-4',
    emerald: 'data-[state=active]:bg-emerald-500 data-[state=active]:text-white',
  };

  return (
    <TabsTrigger
      ref={ref}
      className={cn(
        'inline-flex items-center justify-center whitespace-nowrap px-3 py-1.5 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
        variantStyles[variant],
        className
      )}
      {...props}
    />
  );
});
ShadcnTabsTrigger.displayName = 'ShadcnTabsTrigger';

const ShadcnTabsContent = React.forwardRef<
  React.ElementRef<typeof TabsContent>,
  React.ComponentPropsWithoutRef<typeof TabsContent>
>(({ className, ...props }, ref) => (
  <TabsContent
    ref={ref}
    className={cn(
      'mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2',
      className
    )}
    {...props}
  />
));
ShadcnTabsContent.displayName = 'ShadcnTabsContent';

export {
  ShadcnTabs as Tabs,
  ShadcnTabsList as TabsList,
  ShadcnTabsTrigger as TabsTrigger,
  ShadcnTabsContent as TabsContent,
};
