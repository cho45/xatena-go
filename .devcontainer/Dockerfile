FROM golang:1.24

# plenv, perl-build, および依存パッケージのインストール
RUN apt-get update && \
    apt-get install -y git curl build-essential libssl-dev libbz2-dev libreadline-dev libsqlite3-dev zlib1g-dev libffi-dev libgdbm-dev libncurses5-dev libdb-dev && \
    git clone https://github.com/tokuhirom/plenv.git ~/.plenv && \
    git clone https://github.com/tokuhirom/Perl-Build.git ~/.plenv/plugins/perl-build && \
    echo 'export PLENV_ROOT="$HOME/.plenv"' >> /etc/profile.d/plenv.sh && \
    echo 'export PATH="$PLENV_ROOT/bin:$PATH"' >> /etc/profile.d/plenv.sh && \
    echo 'eval "$(plenv init -)"' >> /etc/profile.d/plenv.sh && \
    chmod +x /etc/profile.d/plenv.sh

# plenvのパスを有効化し、Perl 5.38.2をインストール
ENV PERL_VERSION=5.38.2
RUN bash -c 'source /etc/profile.d/plenv.sh && plenv install $PERL_VERSION && plenv global $PERL_VERSION && plenv rehash'

ENV LANG=C
RUN plenv install-cpanm

RUN cpanm -f UNIVERSAL::require Exporter::Lite Text::MicroTemplate List::MoreUtils URI::Escape  HTML::Entities
