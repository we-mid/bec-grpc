# How to add plugin to RabbitMQ docker image?
# https://stackoverflow.com/questions/52819237/how-to-add-plugin-to-rabbitmq-docker-image
# https://blog.rabbitmq.com/posts/2015/04/scheduling-messages-with-rabbitmq/
# FROM rabbitmq:3.9-management
FROM rabbitmq:3-management

# COPY rabbitmq.conf /etc/rabbitmq/rabbitmq.conf

RUN apt-get -o Acquire::Check-Date=false update && apt-get install -y curl

# RUN curl -L https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/3.9.0/rabbitmq_delayed_message_exchange-3.9.0.ez\
#   > $RABBITMQ_HOME/plugins/rabbitmq_delayed_message_exchange-3.9.0.ez
RUN curl -L https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.12.0/rabbitmq_delayed_message_exchange-3.12.0.ez \
  > $RABBITMQ_HOME/plugins/rabbitmq_delayed_message_exchange-3.12.0.ez

# RUN chown rabbitmq:rabbitmq $RABBITMQ_HOME/plugins/rabbitmq_delayed_message_exchange-3.9.0.ez
RUN chown rabbitmq:rabbitmq $RABBITMQ_HOME/plugins/rabbitmq_delayed_message_exchange-3.12.0.ez

RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange
