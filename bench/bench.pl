#!/usr/bin/env perl

use lib './Text-Xatena/lib';
use Text::Xatena;
use Time::HiRes qw(gettimeofday tv_interval);

my $input = do { local $/; <> };
my $xatena = Text::Xatena->new;

my $t0 = [gettimeofday];
my $html;
for (1..1000) {
    $html = $xatena->format($input);
}
my $elapsed = tv_interval($t0);

print STDERR "pl parse+format (1000x): $elapsed ms\n";
# print $html; # 出力は不要なら省略
