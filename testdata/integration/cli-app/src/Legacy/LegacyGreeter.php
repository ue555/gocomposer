<?php

class LegacyGreeter
{
    public static function suffix(): string
    {
        return \Monolog\Level::Info->getName() . ' via classmap';
    }
}
