# Use Node 20 or 22
FROM node:22

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json first
COPY package.json package-lock.json ./

# Install dependencies (including devDependencies)
RUN npm install

# Copy the rest of the source code
COPY . .

# Expose Vite port
EXPOSE 5173

# Run dev server
CMD ["npm", "run", "dev"]
