# FROM ubuntu:20.04

# # Install dependencies
# RUN apt-get update && apt-get install -y curl git build-essential

# # Install Foundry
# RUN curl -L https://foundry.paradigm.xyz | bash && \
#     /root/.foundry/bin/foundryup

# # Add Foundry to PATH
# ENV PATH="/root/.foundry/bin:${PATH}"

# # Set the working directory
# WORKDIR /project

# # Copy the project files into the container
# COPY ./test-erc20-project /project

# # Set the command to deploy the contract
# CMD ["sh", "-c", "forge script ./script/IERC20Deployer.s.sol:IERC20Deployer --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast -vvvvv"]